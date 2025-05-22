# https://github.com/charmbracelet/vhs
# Requires access to account 308 to execute the examples, switch context if needed first.

Output safedns-example.gif

Require ans

Set FontFamily "VictorMono Nerd Font"
Set FontSize 13
Set TypingSpeed 0.1
Set Theme "catppuccin-macchiato"
Set WindowBar Colorful
Set Margin 25
Set MarginFill "#202e52"
Set Padding 5
Set Width 1200
Set Height 500

Type@20ms "# https://developers.ukfast.io/api-guide#filtering"
Sleep 0.1
Enter
Enter
Type "ans safedns zone list --filter 'name=ans.*uk' --filter 'name:nlk=*.net.uk'"
Sleep 0.3
Enter
Sleep 2

Type "ans safedns record list ans.co.uk --filter 'name=ans.co.uk' --filter 'type:in=A,NS'"
Sleep 0.3
Enter
Sleep 2

Type "ans safedns record create ans.co.uk --name cli-example.ans.co.uk --type A --content 1.2.3.4"
Sleep 0.3
Enter
Sleep 2

Type@50ms "RECORD_ID=$(ans safedns record list ans.co.uk --filter 'name:eq=cli-example.ans.co.uk' --output value --property ID)"
Sleep 0.3
Enter
Sleep 1
Type "echo $RECORD_ID"
Sleep 0.3
Enter
Sleep 1
Type "ans safedns record delete ans.co.uk ${RECORD_ID}"
Sleep 0.3
Enter
Sleep 2
