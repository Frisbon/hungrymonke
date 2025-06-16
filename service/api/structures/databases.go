package structures

import "sync"

// implemento il mutex come richiesto dal professore.
var DBMutex = &sync.RWMutex{} // read write mutex

/*

	Ora dovrò fare una lock all'inizio di ogni sezione che modifica (WRITE) i dati
	sul DB e unlock subito dopo (con defer, in caso di errori)

	scs.DBMutex.Lock() -- Lock per scrivere
	defer scs.DBMutex.Unlock() -- sblocca alla fine

	Ogni volta che leggo i dati uso una RLock perchè più efficiente e non blocca gli altri lettori.

*/

var GenericDB = make(map[string]struct{}) //  Database/Set con tutti gli ID generici, in modo che possa evitare di avere doppioni

var UserDB = make(map[string]*User)       //  nickname : struct utente.   //  Hint: siccome è una mappa, posso avere solo una key con un certo nickname => se esiste record per key, nickname occupato.
var PrivateDB = make(map[string]*Private) //  DB con tutte le chat "1v1"  sarà una mappa IDgenerico_conversazione : Conversazione Privata
var GroupDB = make(map[string]*Group)
var MsgDB = make(map[string]*Message)
var ConvoDB = make(map[string]*ConversationELT)

var UserConvosDB = make(map[string]Conversations) //   nickname : lista conversazioni di quell'utente.
