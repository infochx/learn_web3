// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract MyERC20 {
    // 任务：参考 openzeppelin-contracts/contracts/token/ERC20/IERC20.sol实现一个简单的 ERC20 代币合约。要求：
    // 1. 合约包含以下标准 ERC20 功能：
    // 2. balanceOf：查询账户余额。
    // 3. transfer：转账。
    // 4. approve 和 transferFrom：授权和代扣转账。
    // 5. 使用 event 记录转账和授权操作。
    // 6. 提供 mint 函数，允许合约所有者增发代币。
    // 提示：
    // - 使用 mapping 存储账户余额和授权信息。
    // - 使用 event 定义 Transfer 和 Approval 事件。
    // - 部署到sepolia 测试网，导入到自己的钱包

    mapping (address account => uint256 value) balances;

    mapping(address account => mapping(address spender => uint256)) public allowances;

    // 代币总供应量
    uint256 public  totalSupply;
    // 精度.表示代币可以分割到的小数位数。许多代币选择18为其小数值，因为这是 Ether(ETH) 使用的小数位数
    uint8 public decimals; 
    // 代币名称
    string public name;
    // 代币简称
    string public symbol;
    // 合约拥有者
    address public owner;


    // 事件定义
    event Mint(address indexed _from, address indexed _to, uint256 _value);
    // 在代币被转移时触发。
    event Transfer(address indexed _from, address indexed _to, uint256 _value);
    //在调用 approve 方法时触发。
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);

    // 自定义错误
    error ERC20InvalidSender(address sender);
    error ERC20InvalidReceiver(address receiver);
    error ERC20InvalidSpender(address spender);
    error OwnableUnauthorizedAccount(address sender);
    error ERC20InsufficientBalance(address from, uint256 fromBalance, uint256 value);
    error ERC20InsufficientAllowance(address spender, uint256 currentAllowance, uint256 value);

    // 自定义修改器
    // 只有使用者可以调用
    modifier onlyOwner() {
        if (owner != msg.sender) {
            revert OwnableUnauthorizedAccount(msg.sender);
        }
        _;
    }

    // 构造方法
    constructor() {
        owner = msg.sender;
        name = "MyToken"; 
        symbol = "MTK"; 
        decimals = 18; 
        totalSupply = 100000000 * 10 ** uint256(decimals);
        balances[msg.sender] = totalSupply;  
    }

    // 铸造代币，只有合约所有者有权增发代币。
    function mint(address account, uint256 value) public onlyOwner {
        totalSupply += value;
        if (account != address(0)) {
            balances[account] += value;
        }
        emit Mint(msg.sender,account,value);
    }

    //  2. balanceOf：查询账户余额。
    function balanceOf(address account) public view returns (uint256) {
        return balances[account];
    }

    // 3. transfer：转账。
    function transfer(address _to,uint256 _value) public returns (bool) {
        if(_to == address(0)){
            revert ERC20InvalidReceiver(_to);
        }
        address sender = msg.sender;
        if(balances[sender] < _value){
            revert ERC20InsufficientBalance(sender, balances[sender], _value);
        }
        balances[sender] -= _value;
        balances[_to] += _value;
        emit Transfer(sender, _to, _value);
        return true;
    }


    // 4. approve 和 transferFrom：授权和代扣转账。
    function approve(address spender, uint256 value) public returns (bool) {
        address _owner = msg.sender;
        if (spender == address(0)) {
            revert ERC20InvalidSpender(address(0));
        }
        allowances[_owner][spender] += value;
        emit Approval(_owner, spender, allowances[_owner][spender]);
        return true;
    }

    function transferFrom(address from, address to, uint256 value) public virtual returns (bool) {
        // 地址校验
        if (from == address(0)) {
            revert ERC20InvalidSender(address(0));
        }
        if (to == address(0)) {
            revert ERC20InvalidReceiver(address(0));
        }

        // 授权信息校验及更新
        address spender = msg.sender;
        uint256 currentAllowance = allowances[from][spender];
        if (currentAllowance < type(uint256).max) {
            if (currentAllowance < value) {
                revert ERC20InsufficientAllowance(spender, currentAllowance, value);
            }
            unchecked {
                // 更新授权金额
                allowances[from][spender] = currentAllowance - value;
            }
        }

        // 完成代币转移
        balances[from] -=value;
        balances[to] += value;
        emit Transfer(from,to,value);
        return true;
    }

    // 返回 _spender 仍然被允许从 _owner 提取的代币数量
    function allowance(address _owner, address _spender) public view returns (uint256 remaining){
        return allowances[_owner][_spender];
    }
}
