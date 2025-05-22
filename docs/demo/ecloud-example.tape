# https://github.com/charmbracelet/vhs
# Note, check VPC ID below and correct for account used.

Output ecloud-example.gif

Require ans

Set FontFamily "VictorMono Nerd Font"
Set FontSize 16
Set TypingSpeed 0.1
Set Theme "catppuccin-macchiato"
Set WindowBar Colorful
Set Padding 5
Set Width 1200
Set Height 500

Type "ans ecloud instance list"
Sleep 0.3
Enter
Sleep 2

Type "ans ecloud network list"
Sleep 0.3
Enter
Sleep 2

Type "ans ecloud image list --filter 'name=AlmaLinux*'"
Sleep 0.3
Enter
Sleep 2

Type@50ms "ans ecloud instance create \"
Enter
Type@50ms "  --vpc vpc-e6613504 \"
Enter
Type@50ms "  --network net-64169be9 \"
Enter
Type@50ms "  --image img-19cb94e5 \"
Enter
Type@50ms "  --volume 40 \"
Enter
Type@50ms "  --vcpu-cores-per-socket 1 \"
Enter
Type@50ms "  --ram 1024 \"
Enter
Type@50ms "  --name 'CLI Hello World'"
Sleep 0.3
Enter
Sleep 4

Type@50ms "INSTANCE_ID=$(ans ecloud instance list --filter 'name=CLI*' --output value --property ID)"
Sleep 0.3
Enter
Sleep 2

Type@50ms "ans ecloud instance show $INSTANCE_ID --output json | jq -r '.[0].sync'"
Sleep 0.3
Enter
Sleep 2

Hide
Sleep 360
Show

Type "# ... time passes ..."
Enter

Type@50ms "ans ecloud instance show $INSTANCE_ID --output json | jq -r '.[0].sync'"
Sleep 0.3
Enter
Wait
Type "ans ecloud instance delete $INSTANCE_ID"
Sleep 0.3
Enter
Wait
Type@50ms "ans ecloud instance show $INSTANCE_ID --output json | jq -r '.[0].sync'"
Sleep 0.3
Enter
Wait
Sleep 5
