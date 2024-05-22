要编写一个启用 Swap 的脚本，并确保在需要时重新启用 Swap，可以按照以下步骤进行：

1. **创建 Swap 文件或分区**：如果之前删除了 Swap 文件或分区，需要重新创建。
2. **编辑 `/etc/fstab` 文件**：重新添加 Swap 条目。
3. **启用 Swap**：使用 `swapon` 命令启用 Swap。

假设之前是使用 Swap 文件，可以使用以下脚本：

```bash
#!/bin/bash

# 创建一个 2GB 的 Swap 文件（根据需要调整大小）
sudo fallocate -l 2G /swapfile

# 设置交换文件权限
sudo chmod 600 /swapfile

# 将文件格式化为 Swap
sudo mkswap /swapfile

# 启用 Swap 文件
sudo swapon /swapfile

# 确认 Swap 是否启用
sudo swapon --show

# 将 Swap 文件信息添加到 /etc/fstab 以便系统重启时自动启用
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab

# 确认 /etc/fstab 文件是否更新
cat /etc/fstab
```

### 解释

1. **创建 Swap 文件**：
   ```bash
   sudo fallocate -l 2G /swapfile
   ```
   - 创建一个大小为 2GB 的空文件 `/swapfile`。
   - 你可以根据需要调整文件大小，比如 `1G`、`4G` 等。

2. **设置交换文件权限**：
   ```bash
   sudo chmod 600 /swapfile
   ```
   - 设置 `/swapfile` 的权限为 600，以确保只有 root 用户可以读取和写入该文件。

3. **将文件格式化为 Swap**：
   ```bash
   sudo mkswap /swapfile
   ```
   - 使用 `mkswap` 命令将文件格式化为 Swap 文件。

4. **启用 Swap 文件**：
   ```bash
   sudo swapon /swapfile
   ```
   - 启用新的 Swap 文件。

5. **确认 Swap 是否启用**：
   ```bash
   sudo swapon --show
   ```
   - 使用 `swapon --show` 命令查看当前启用的 Swap 设备。

6. **将 Swap 文件信息添加到 `/etc/fstab`**：
   ```bash
   echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
   ```
   - 将 Swap 文件的信息添加到 `/etc/fstab` 文件中，以便系统在重启时自动启用该 Swap 文件。

7. **确认 `/etc/fstab` 文件是否更新**：
   ```bash
   cat /etc/fstab
   ```
   - 查看 `/etc/fstab` 文件，确认 Swap 文件信息已添加。

### 使用方法

将上述脚本保存为一个文件，例如 `enable_swap.sh`，然后赋予执行权限并运行该脚本：

```bash
chmod +x enable_swap.sh
./enable_swap.sh
```

通过这个脚本，您可以在需要时轻松启用 Swap，并确保在系统重启后自动启用 Swap。