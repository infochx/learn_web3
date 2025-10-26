// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./interfaces/IAuctionFactory.sol";

contract Auction is ReentrancyGuard {
    // 拍卖工厂合约地址
    address public factory;
    // NFT拍卖发起人地址
    address public seller;
    // NFT合约地址
    address public nftContract;
    // NFT令牌ID
    uint256 public nftTokenId;
    // 拍卖开始时间
    uint256 public startTime;
    // 拍卖结束时间
    uint256 public endTime;
    // 初始拍卖价格
    uint256 public startPrice;
    // 最高出价
    uint256 public highestBid;
    // 最高出价者地址
    address public highestBidder;
    // 货币类型 0x00: ETH 其它：LINK
    address public tokenAddress;
    // 货币真实数量
    uint256 public tokenValue = 0;
    // 拍卖是否结束
    bool public isEnded;
    // 出价记录
    mapping(address => uint256) public bidsData;

    modifier onlyFactory() {
        require(msg.sender == factory, "Only factory can call this function");
        _;
    }

    modifier onlySeller() {
        require(msg.sender == seller, "Only seller can call this function");
        _;
    }

    modifier timeToBid() {
        require(
            block.timestamp >= startTime && block.timestamp <= endTime,
            "Auction not start or end"
        );
        _;
    }

    fallback() external payable {}

    receive() external payable {}

    constructor() {
        factory = msg.sender;
    }

    function initialize(
        address _seller,
        address _nftContract,
        uint256 _nftTokenId,
        uint256 _startPrice,
        uint256 _startTime,
        uint256 _duration
    ) external onlyFactory {
        seller = _seller;
        nftContract = _nftContract;
        nftTokenId = _nftTokenId;
        startPrice = _startPrice;
        startTime = _startTime;
        endTime = _startTime + _duration;
    }

    function placeBidWithETH() external payable nonReentrant timeToBid {
        require(!isEnded, "Auction already ended");
        uint256 amount = IAuctionFactory(factory).formatEthToUsdtPrice(
            msg.value
        );

        require(
            amount > highestBid && amount >= startPrice,
            "Bid not high enough"
        );

        // 退还之前的最高出价者的ETH
        returnToken();

        // 更新最高出价者
        highestBidder = msg.sender;
        // 最高出价对应USDT，用于比较
        highestBid = amount;
        tokenAddress = address(0);
        // 最高出价ETH，用于退款
        tokenValue = msg.value;
        // 更新出价记录
        bidsData[msg.sender] = amount;
    }

    function placeBidWithERC20(
        address bidTokenAddress,
        uint256 value
    ) external nonReentrant timeToBid {
        require(!isEnded, "Auction already ended");
        
        require(factory != address(0), "Invalid factory address");
        uint256 amount = IAuctionFactory(factory).formatLinkToUsdtPrice(value);
        require(
            amount > highestBid && amount >= startPrice,
            "Bid not high enough"
        );
        bool success = IERC20(bidTokenAddress).transferFrom(msg.sender, address(this), value);
        require(success, "Transfer failed");

        // 退还之前的最高出价者的ERC20
        returnToken();

        // 更新最高出价者
        highestBidder = msg.sender;
        // 最高出价对应USDT，用于比较
        highestBid = amount;
        tokenAddress = bidTokenAddress;
        // 最高出价ERC20，用于退款
        tokenValue = value;
        // 更新出价记录
        bidsData[msg.sender] = amount;
    }

    // 退还最高出价者的ETH或ERC20
    function returnToken() private{
        // 如果有最高出价者
        if(highestBidder != address(0)){
            // 如果是ETH
            if(tokenAddress == address(0)){
                // 退还最高出价者的ETH
                (bool success,bytes memory data ) = highestBidder.call{value: tokenValue}("");
                require(success && (data.length == 0), "Transfer failed");
            }else {
                // 退还最高出价者的ERC20
                require(IERC20(tokenAddress).transfer(highestBidder, tokenValue), "Transfer failed");
            }
        }
    }

    function cancelAuction() external onlySeller{
        isEnded = true;
        IERC721(nftContract).transferFrom(address(this), seller, nftTokenId);
        // 退还所有出价者的ETH或ERC20
        returnToken();

        // 通知工厂拍卖结束
        IAuctionFactory(factory).auctionEnd(address(this));
    }

    function endAuction() external onlySeller{
        // 1. 检查拍卖是否结束
        require(!isEnded, "Auction already ended");
        // 2. 检查是否到了结束时间
        require(block.timestamp >= endTime, "Auction not end");


        // 3. 是否有最高出价者
        if(highestBidder != address(0)){
            // 3.1 有最高出价者，转移NFT给最高出价者
            IERC721(nftContract).transferFrom(address(this), highestBidder, nftTokenId);

            // 3.2 货币类型是为ETH
            if(tokenAddress == address(0)){
                // 将ETH转账给卖家
                (bool success,bytes memory data ) = payable(seller).call{value: tokenValue}("");
                require(success && (data.length == 0), "Transfer failed");
            }else {
                // 将ERC20转账给卖家
                require(IERC20(tokenAddress).transfer(seller, tokenValue), "Transfer failed");
            }

        }else{
            // 4. 没有最高出价者，返回NFT给卖家
            IERC721(nftContract).transferFrom(address(this), seller, nftTokenId);
        }
        isEnded = true;
        // 5. 通知工厂拍卖结束
        IAuctionFactory(factory).auctionEnd(address(this));
    }

    /**
     * @notice 处理ERC721接收事件。
     * 如果没有实现这个方法，使用 safeTransferFrom 将NFT转移到拍卖合约时会失败，导致拍卖无法创建。
     */
    function onERC721Received(
        address,
        address,
        uint256,
        bytes calldata
    ) external pure returns (bytes4) {
        return this.onERC721Received.selector;
    }
}
