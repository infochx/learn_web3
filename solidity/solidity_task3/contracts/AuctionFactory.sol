// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./Auction.sol";
import "./interfaces/IAuctionFactory.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "hardhat/console.sol";

// 拍卖市场，可以创建拍卖
contract AuctionFactory is IAuctionFactory {
    // 拍卖市场的所有者
    address private marketOwner;
    // 所有的拍卖
    address[] public auctions;
    // 下一个拍卖的ID
    uint256 private _nextAuctionId = 1;
    // 拍卖ID => 拍卖合约地址
    mapping(uint256 auctionId => address auctionContract) private _auctionData;

    // 货币地址 => 价格Feed地址
    mapping(address => AggregatorV3Interface) private _priceFeeds;

    address linkAddress = 0x779877A7B0D9E8603169DdbD7836e478b4624789;
    address ethAddress = 0x0000000000000000000000000000000000000000;

    event AuctionCreate(address indexed auctionAddress);
    bool private _isHardHat = false;

    constructor(bool isHardHat, address _linkAddress) {
        marketOwner = msg.sender;
        if (_linkAddress != address(0)) {
            linkAddress = _linkAddress;
        }

        //linkAddress 如果使用本地地址（如 Hardhat 本地网络部署的 MockLINK 合约地址），
        //会导致 价格喂价（AggregatorV3Interface）调用失败，因为你绑定的
        //0xc59E3633BAAC79493d908e63626716e204A45EdF
        //和 0x694AA1769357215DE4FAC081bf1f309aDC325306
        //大概率是 Sepolia 测试网的 Chainlink 价格喂价地址，在本地网络中这些地址是无效的（没有部署对应的合约）
        _priceFeeds[linkAddress] = AggregatorV3Interface(
            0xc59E3633BAAC79493d908e63626716e204A45EdF
        );
        _priceFeeds[ethAddress] = AggregatorV3Interface(
            0x694AA1769357215DE4FAC081bf1f309aDC325306
        );
        _isHardHat = isHardHat;
    }

    function createAuction(
        address _nftContract,
        uint56 _nftTokenId,
        uint256 _startPrice,
        uint256 _duration
    ) public returns (uint256) {
        Auction auction = new Auction();
        uint256 _startPriceUsdt = formatEthToUsdtPrice(_startPrice);

        auction.initialize(
            msg.sender,
            _nftContract,
            _nftTokenId,
            _startPriceUsdt,
            block.timestamp,
            _duration
        );
        // 将NFT转让给拍卖合约
        IERC721(_nftContract).transferFrom(
            msg.sender,
            address(auction),
            _nftTokenId
        );

        address auctionAddress = address(auction);
        // 存储拍卖合约地址
        auctions.push(auctionAddress);
        // 存储拍卖ID => 拍卖合约地址
        _auctionData[_nextAuctionId] = auctionAddress;

        emit AuctionCreate(auctionAddress);
        // 返回拍卖ID
        return _nextAuctionId++;
    }

    // 获取拍卖合约数量
    function getAuctionCount() public view returns (uint256) {
        return auctions.length;
    }

    // 获取拍卖Id对应的合约地址
    function getAuctionAddress(
        uint256 auctionId
    ) public view returns (address) {
        return _auctionData[auctionId];
    }

    // 结束某个拍卖
    function auctionEnd(address auctionAddress) external override {
        uint256 length = auctions.length;
        for (uint256 i = 0; i < length; i++) {
            if (auctions[i] == auctionAddress) {
                auctions[i] = auctions[length - 1];
                auctions.pop();
                break;
            }
        }
    }

    // 获取LINK价格（单位：USDT）
    function formatLinkToUsdtPrice(
        uint256 amount
    ) public view returns (uint256) {
        if (_isHardHat) {
            return amount * uint256(10 ** 8);
        }

        // 调用 Chainlink 价格喂价合约获取 LINK 价格（单位：USDT）
        (, int256 price, , , ) = _priceFeeds[linkAddress].latestRoundData();
        return amount * uint256(price);
    }

    // 获取ETH价格（单位：USDT）
    function formatEthToUsdtPrice(
        uint256 amount
    ) public view returns (uint256) {
        if (_isHardHat) {
            return amount * uint256(10 ** 8);
        }

        // 调用 Chainlink 价格喂价合约获取 ETH 价格（单位：USDT）
        (, int256 price, , , ) = _priceFeeds[ethAddress].latestRoundData();
        return amount * uint256(price);
    }
}
