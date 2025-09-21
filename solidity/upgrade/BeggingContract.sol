// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

// 任务目标

// 1. 使用 Solidity 编写一个合约，允许用户向合约地址发送以太币。
// 2. 记录每个捐赠者的地址和捐赠金额。
// 3. 允许合约所有者提取所有捐赠的资金。

// 任务步骤

// 1. 编写合约
//    - 创建一个名为 BeggingContract 的合约。
//    - 合约应包含以下功能：
//    - 一个 mapping 来记录每个捐赠者的捐赠金额。
//    - 一个 donate 函数，允许用户向合约发送以太币，并记录捐赠信息。
//    - 一个 withdraw 函数，允许合约所有者提取所有资金。
//    - 一个 getDonation 函数，允许查询某个地址的捐赠金额。
//    - 使用 payable 修饰符和 address.transfer 实现支付和提款。
// 2. 部署合约
//    - 在 Remix IDE 中编译合约。
//    - 部署合约到 Goerli 或 Sepolia 测试网。
// 3. 测试合约
//    - 使用 MetaMask 向合约发送以太币，测试 donate 功能。
//    - 调用 withdraw 函数，测试合约所有者是否可以提取资金。
//    - 调用 getDonation 函数，查询某个地址的捐赠金额。

// 任务要求

// 1. 合约代码：
//    - 使用 mapping 记录捐赠者的地址和金额。
//    - 使用 payable 修饰符实现 donate 和 withdraw 函数。
//    - 使用 onlyOwner 修饰符限制 withdraw 函数只能由合约所有者调用。
// 2. 测试网部署：
//    - 合约必须部署到 Goerli 或 Sepolia 测试网。
// 3. 功能测试：
//    - 确保 donate、withdraw 和 getDonation 函数正常工作。

// 提交内容

// 1. 合约代码：提交 Solidity 合约文件（如 BeggingContract.sol）。
// 2. 合约地址：提交部署到测试网的合约地址。
// 3. 测试截图：提交在 Remix 或 Etherscan 上测试合约的截图。

// 额外挑战（可选）

// 1. 捐赠事件：添加 Donation 事件，记录每次捐赠的地址和金额。
// 2. 捐赠排行榜：实现一个功能，显示捐赠金额最多的前 3 个地址。
// 3. 时间限制：添加一个时间限制，只有在特定时间段内才能捐赠。
contract BeggingContract {
// 任务目标

// 1. 使用 Solidity 编写一个合约，允许用户向合约地址发送以太币。
// 2. 记录每个捐赠者的地址和捐赠金额。
// 3. 允许合约所有者提取所有捐赠的资金。

// 任务步骤

// 1. 编写合约
//    - 创建一个名为 BeggingContract 的合约。
//    - 合约应包含以下功能：
//    - 一个 mapping 来记录每个捐赠者的捐赠金额。
mapping (address account => uint256 value) public _donors;

address public _owner;

uint256[3] public _topThree;
address[3] public _topDonors;

uint256 deploymentTimestamp;

uint256 lockTime;

constructor(uint256 _lockTime) {
    _owner = msg.sender;
    deploymentTimestamp = block.timestamp;// 秒
    lockTime = _lockTime;

}


// 自定义error
error InsufficientValue(address donor,uint256 value);
error OnlyOwner(address owner,address callUser);
error Withdraw(address owner,uint256 value);
// 自定义事件
event DonatedEvent(address donor,uint256 value);
event WithdrawEvent(address owner,uint256 value);

// 自定义modifer
modifier onlyOwner(){
    if(msg.sender != _owner){
        revert OnlyOwner(_owner,msg.sender);
    }
    _;
}


modifier windowClosed() {
    require(block.timestamp <= deploymentTimestamp + lockTime, "window is closed");
    _;
}

//    - 一个 donate 函数，允许用户向合约发送以太币，并记录捐赠信息。
function donate() public payable windowClosed returns (bool){
    uint256 _value = msg.value;
    address _sender = msg.sender;
    if(_value == 0){
        revert InsufficientValue(_sender,_value);
    }
    _donors[_sender] += _value;
    // 更新排行榜
    _updateTop(_sender,_donors[_sender]);
    emit DonatedEvent(_sender,_value);
    return true;
}

function _updateTop(address _donor,uint256 _value) internal {
    bool flag = false;
    uint256 value = _value;
    address donor = _donor;
    uint256 currentValue ;
    for (uint256 index; index<3; index++){

        currentValue = _topThree[index];
        if(currentValue < value){
            uint256 value_ = _topThree[index]; 
            address donor_ = _topDonors[index];
            _topThree[index] = value;
            _topDonors[index] = donor;
            value = value_;
            donor = donor_;
            
            flag = true;
            continue ;
        }
        if(flag){
            uint256 value_ = _topThree[index]; 
            address donor_ = _topDonors[index];
            _topThree[index] = value;
            _topDonors[index] = donor;
            value = value_;
            donor = donor_;
        }
    }
}

//    - 一个 withdraw 函数，允许合约所有者提取所有资金。
function withdraw() public onlyOwner returns(bool){
    uint256 _balance = address(this).balance;
    if (_balance == 0){
        revert InsufficientValue(_owner,_balance);
    }
    (bool success, ) = payable(_owner).call{value: _balance}("");
    if(!success){
        revert Withdraw(_owner,_balance);
    }
    emit WithdrawEvent(_owner,_balance);
    return true;
}

//    - 一个 getDonation 函数，允许查询某个地址的捐赠金额。
function getDonation(address _donor) public view returns(uint256){
    return _donors[_donor];
}


}
