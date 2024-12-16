# Tắt swap
sudo swapoff -a

# Tắt Address Space Randomisation
echo 0 > /proc/sys/kernel/randomize_va_space

# Tắt Transparent Hugepages
echo never > /sys/kernel/mm/transparent_hugepage/enabled
echo never > /sys/kernel/mm/transparent_hugepage/defrag
echo 0 > /sys/kernel/mm/transparent_hugepage/khugepaged/defrag

# Cấu hình core dump
echo core > /proc/sys/kernel/core_pattern
isolate --cg --cleanup 
isolate --cleanup
arr=$(ls /var/local/lib/isolate/)
for i in "${arr[@]}"; do
echo "isolate --cg --cleanup  --box-id=$i"
isolate --cg --cleanup  --box-id=$i
done