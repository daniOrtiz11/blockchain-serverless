package main

//AWS CONSTANTS:
const fileUploadError2 = "file uploaded to %s\n"
const responseParam = "RequestResponse"
const logParam = "Tail"
const localfileblc string = "blc.json"
const localfilebank string = "bank.json"
const okC = "ok"
const koC = "ko"
const emptyC = "empty"
const blockchainstr = "blockchain"
const initParam = "{\"newnode\": \""

//P2P-BLC CONSTANTS:
const ipfs = "/ipfs/%s"
const ipfs2 = "/ipfs/"
const ipv4 = "/ip4"
const ipv4_2 = "/ip4/"
const tcp = "tcp"
const tcp2 = "/tcp/"
const iplocalhost = "/ip4/127.0.0.1/tcp/%d"
const localhost = "127.0.0.1"
const p2p = "/p2p/1.0.0"
const initSource = "sourceGenesis"
const initTarget = "targetGenesis"
const initAmount = -2
const createAmountName = -1
const createAmountPriv = 0
const defaultPort = 10000
const initAmountAccount = 100

//ERROR CONSTANTS:
const exceptionJSON = "Exception at Json library: "
const exceptionReader = "Exception in reader: "
const badFormatOption = "Error: Please provide a option"
const badFormatNumber = "Error: Please provide a unsigned number"
const badFormatArgument = "Please provide a port to bind on with -l"
const errorSetP2P = "Error setting P2P node"
const errorDeleteP2P = "Error deleting P2P node"
const errorSetLog = "Error updating log file"
const errorNetwork = "Error trying to get your ip"
const errorDirRepeat = "Error establishing the connection. It is possible that the key, ip or port are repeated."
const errorLogin = "Login error. Check your private key."
const errorCreateAccount = "Error in account creation. Name in use, choose another one."
const errorRespAws = "Error in response type parsing AWS."
const awsError = "Error getting an answer from AWS."
const fileOpenError = "Error opening the file %q, %v"
const fileUploadError = "Error uploading the file %v"

//MENU CONSTANTS:
const helpRunning = "Try providing a different port with option -l"
const defaultPortStr = "Running in default port 10000"
const greenColorCMD = "\x1b[32m%s\x1b[0m> "
const cmdInitialNode = "Initial node"
const cmdInitialNode2 = "Waiting for connections..."
const cmdClientNode = "Opening stream"
const flagL = "l"
const flagSeed = "seed"
const startingSetP2P = "Starting the node setting..."
const startedSetP2P = "Set successfully!"
const koPingP2P = "Your node is currently disconnected"
const reconnectingP2P = "Reconnecting..."
const nameAccountStr = "Name: "
const coinsAccountStr = "Coins: "
const publicAccountStr = "Public Key: "
const privateAccountStr = "Private Key: "

//ACTIONS CONSTANTS:
const optionsTitle = "Options:"
const options1Title = "1. View Blockchain state"
const options2Title = "2. View your account state"
const options2Title1 = "Current account state:"
const options1Title1 = "Current Blockchain state:"
const options3Title = "3. New transaction"
const options3Title11 = "Enter the number of coins to transfer:"
const options3Title12 = "You don't have enough coins in your account. Try less!"
const options3Title21 = "Type the user id:"
const options3Title22 = "Wrong id! Try again!"
const options3Title3 = "Sending the transfer to "
const options4Title = "4. View bank state"
const options4Ttile1 = "Current bank state:"
const closeOptionMenu = "5. Close connection"
const login1 = "1. Log in"
const login1Title1 = "Enter your private key:"
const welcomeLogged = "Welcome "
const login2 = "2. Create account"
const login2Title1 = "Type your new name:"
const login2Title2 = "The account has been created!"
const login2Title3 = "Your account is:"
const closeOptionLogin = "3. Close connection"
const endMessage = "Closing connection, bye!"
const logDelimeter = "|"
const logColor = "#FFFFFF"
const localfilelog = "blockchain-serverless.log"
const extnamelogs = ".accounts"
const folder1logs = "bank/"
const folder2logs = "blockchain/"
