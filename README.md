# Bitmex-referral-analyzer

This application can be used in 2 ways.

You can put all your Bitmex .csv history files in the csv folder and run main.go

The other way is to run a Telegram bot and analyze your referral earnings.
Put in the API keys for telegram and Bitmex in the config folder.
###### Minimum Recommended Specifications

- Go 1.10 or 1.11
* Linux


  Installation instructions can be found here: https://golang.org/doc/install.
  It is recommended to add `$GOPATH/bin` to your `PATH` at this point.

###### setup
``cd ~/go/src/github.com/``

``git clone git@gitlab.com:romanornr/Bitmex-referral-analyzer``

``cd Bitmex-referral-analyzer``

``dep ensure`` 


dep is a dependency management tool for Go. It requires Go 1.9 or newer to compile.
https://github.com/golang/dep

##### Method 1: CSV file
Your wallet history csv file should be saved in the same folder as the program.

###### Screenshots
![alt text](https://github.com/romanornr/Bitmex-referral-analyzer/blob/master/screenshots/save-as-csv.png?raw=true)
<br><br>


![alt text](https://github.com/romanornr/Bitmex-referral-analyzer/blob/master/screenshots/screenshot.png?raw=true)

##### Method 2: Telegram commands
week - Show earnings this week

months - Show earnings per month

status - Show total referrals and pending payout

###### Ref shill
Join to get that 10% Shitmex discount: https://www.bitmex.com/register/vhT2qm


