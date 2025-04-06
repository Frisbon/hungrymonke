/*
Questo file servirà per mandare tutte le richieste al back-end.
Così posso accedere alle funzioni qui, nelle altre schermate vue
*/  

import axios from 'axios';

const apiClient = axios.create({
  baseURL: 'http://localhost:8082/api',
  headers: {
    'Content-Type': 'application/json', // implica che riceverà solo JSON dal back-end
  },
});

  
// TODO, fai in modo che anche lo username si salvi per app.vue e convolist

  // Funzione per impostare/rimuovere il token negli header di Axios per login e logout
  function setAuthToken(token) {

    if (token) { // se passo il token, aggiungilo
      apiClient.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      localStorage.setItem('token', token); // salvo il token nella localStorage (cache)

    } else { // altrimenti assumo che già ci sia, lo rimuovo
      delete apiClient.defaults.headers.common['Authorization'];
      localStorage.removeItem('token');
    }

  }
  
  // Recupera il token salvato (nella cache? boh) all'avvio
  const savedToken = localStorage.getItem('token');
  if (savedToken) {
    setAuthToken(savedToken);
  }
  
  export default {


  getConvoInfo(convoID){
    const token = localStorage.getItem('token');
    return apiClient.get(`/utils/getconvoinfo/${convoID}`, {
      headers: { Authorization: `Bearer ${token}` }
    });

  },

  getConversations() {

        /*
        Invio una richiesta per le conversazioni and
          - modifico l'header per aggiungere il token
          - subito dopo vedo se ho ricevuto le convo dal back-end
            - quindi obj del tipo {Username: string, User Conversations: [username: array_convo]}
          - controllo se ho un errore.
        */
        return apiClient.get('/conversations', {
                                  headers: {
                                    'Authorization': `Bearer ${localStorage.getItem('token')}`,
                                    // recupero il token salvato nella localStorage
                                  },
        })
        .then(response => {
                            console.log('Conversations fetched:', response.data);
                            return response.data;
        })
        .catch(error => {
                          console.error('Error fetching conversations:', error);
        });
  },

  getMessages(convoID) {
    const token = localStorage.getItem('token');
    return apiClient.get(`/conversations/${convoID}`, {
      headers: { Authorization: `Bearer ${token}` }
    });
  },

  sendMessage(convoID, message) {
    const token = localStorage.getItem('token');
    return apiClient.post(`/conversations/messages?ID=${convoID}`, 
      { ...message }, 
      { headers: { Authorization: `Bearer ${token}` } }
    ); 
  },

  changeUsername(newName) {
    const token = localStorage.getItem('token');
    return apiClient.put(`/users/me/username`, newName, 
      { headers: { Authorization: `Bearer ${token}` } }
    ).then(response => 
        {
          if (response.data && response.data.error) {
          // Scenario 1: Ricevuto JSON con solo il campo 'error'
          console.error("Errore durante il cambio username:", response.data.error);
          return response.data.error
        } else if (response.data && response.data.message && response.data.user && response.data.new_token) {
          // Scenario 2: Ricevuto JSON con { message, user, new_token }
          console.log("Username cambiato con successo:", response.data.message);
          setAuthToken(response.data.new_token)
          return response.data;
        } else {
          console.warn("Risposta in formato inatteso:", response);
          return response;
        }
    })

    
  
    
  },


  
  /*
  Questa funzione invia la stringa "credentials" (nickname) nel body della richiesta al back-end
  Il back-end converte il body in una stringa e usa il nickname ricevuto per loggare l'utente,
  ritornando un JSON. Quindi response.data = {token: String, user: UserStruct}
  */
  async login(credentials) { 
    const response = await apiClient.post('/login', credentials); //invio post al back-end per fare il login
    const token = response.data.token;

    console.log("Ho appena loggato l'utente ("+response.data.user.username+")")
    console.log("Ecco il response del login in api.js:")
    console.log(response)

    setAuthToken(token);
    return response;
  },

  logout() {
    setAuthToken(null);
  },

};
