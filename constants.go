package main

//AWS CONSTANTS:
const awsError = "err to get response from aws"
const fileOpenError = "failed to open file %q, %v"
const fileUploadError = "failed to upload file, %v"
const fileUploadError2 = "file uploaded to, %s\n"

//TEMP FILE CONSTANTS:
const localfile string = "blc.json"

//STRINGS
const ipfs = "/ipfs/%s"
const ipfs2 = "/ipfs/"
const ipv4 = "/ip4"
const ipv4_2 = "/ip4/"
const tcp = "tcp"
const tcp2 = "/tcp/"
const iplocalhost = "/ip4/127.0.0.1/tcp/%d"

const listSources = "main.go bcfunctions.go constants.go aws.go utils.go handlerserver.go"
const exceptionJSON = "Exception in json library: "
const exceptionReader = "Exception in reader: "
const badFormatOption = "Error: Please provide a option"
const badFormatNumber = "Error: Please provide a number"
const badFormatArgument = "Please provide a port to bind on with -l"
const badFormatMode = "Please provide a correct mode with -mode"
const greenColorCMD = "\x1b[32m%s\x1b[0m> "
const debugstr = "debug"
const prodstr = "prod"
const cmdInitialNode = "Initial node"
const p2p = "/p2p/1.0.0"
const cmdInitialNode2 = "Waiting for connections..."
const cmdClientNode = "opening stream"
const flagL = "l"
const flagD = "d"
const flagSeed = "seed"
const flagMode = "mode"

//ACTIONS
const optionsTitle = "Options:"
const options1Title = "1. View BLC State"
const options2Title = "2. View your account State"
const options1Title1 = "Current Blockchain state:"
const options3Title = "3. New Transaction"
const options3Title11 = "Type the amount to transfer:"
const options3Title12 = "You don't have enougth amount in your account! Try again!"
const options3Title21 = "Type the user id:"
const options3Title22 = "Wrong id! Try again!"
const options3Title3 = "Sending amount to "
const options4Title = "4. View Bank"

const closeOptionMenu = "5. Close connection"
const login1 = "1. Log in"
const login2 = "2. Create account"
const closeOptionLogin = "3. Close connection"

const endMessage = "Closing connection, bye!"
const startingSetP2P = "Starting set node..."
const startedSetP2P = "Setted ok!"
const errorSetP2P = "Error setting P2P node"
const errorDeleteP2P = "Error deleting P2P node"
const koPingP2P = "Your node is currently disconnect"
const reconnectingP2P = "Reconnecting..."
const errorDirRepeat = "Error while starting the connection. Possible repeated key, ip or port."
const maxNodes = 100
