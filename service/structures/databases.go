package structures

var GenericDB = make(map[string]struct{}) // Database/Set con tutti gli ID generici, in modo che possa evitare di avere doppioni

var UserDB = make(map[string]*User)       // nickname : struct utente.   // Hint: siccome è una mappa, posso avere solo una key con un certo nickname => se esiste record per key, nickname occupato.
var PrivateDB = make(map[string]*Private) // DB con tutte le chat "1v1"  sarà una mappa IDgenerico_conversazione : Conversazione Privata
var GroupDB = make(map[string]*Group)
var MsgDB = make(map[string]*Message)
var ConvoDB = make(map[string]*ConversationELT)

var UserConvosDB = make(map[string]Conversations) //  nickname : lista conversazioni di quell'utente.
