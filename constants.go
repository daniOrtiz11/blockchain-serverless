package main

//AWS CONSTANTS:
const awsError = "err to get response from aws"
const fileOpenError = "failed to open file %q, %v"
const fileUploadError = "failed to upload file, %v"
const fileUploadError2 = "file uploaded to, %s\n"

//TEMP FILE CONSTANTS:
const localfile string = "blc.json"

//ENV
var debug = true

//STRINGS
const ipfs = "/ipfs/%s"
const ipfs2 = "/ipfs/"
const ipv4 = "/ip4"
const ipv4_2 = "/ip4/"
const tcp = "tcp"
const tcp2 = "/tcp/"
const iplocalhost = "/ip4/127.0.0.1/tcp/%d"

const listSources = "main.go bcfunctions.go constants.go aws.go utils.go handlerserver.go"
const debugcmd = "Now run \"go run %s -l %d -d %s -mode debug\" on a different terminal\n"
const prodcmd = "Now run \"blockchain-serverless -l %d -d %s\" on a different terminal\n"
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
const cmdInitialNode2 = "Listening for connections"
const cmdClientNode = "opening stream"
const flagL = "l"
const flagD = "d"
const flagSeed = "seed"
const flagMode = "mode"
const optionsTitle = "Options:"
const options1 = "1. View State"
const options1Title = "Current Blockchain state:"
const options2 = "2. Insert new value"
const options2Title = "Insert a new value:"
const options3 = "3. Show command to run"
const options4 = "4. Close connection"
const endMessage = "Closing connection, bye!"
const startingSetP2P = "Starting set new node..."
const startedSetP2P = "Setted ok!"
const errorSetP2P = "Error setting P2P"
