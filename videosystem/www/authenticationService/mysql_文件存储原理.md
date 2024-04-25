# 问题背景：
- MySQL 一行记录是怎么存储的？
- MySQL 的 NULL 值会占用空间吗？
- MySQL 怎么知道 varchar(n) 实际占用数据的大小？
- varchar(n) 中 n 最大取值为多少？
- 行溢出后，MySQL 是怎么处理的？

### 查看数据库文件命令
- show variables like 'datadir';

### 表空间结构
1. 段
2. 区
    - 默认大小 1 MB
    - 
3. 页
    - 单次IO的基本单位
    - 默认大小16K
    - 常见类型 数据页 undo 溢出页
4. 行