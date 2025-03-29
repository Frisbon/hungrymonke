import axios from 'axios';

const apiClient = axios.create({
    baseURL: 'http://localhost:8082/api',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  
  // Funzione per impostare il token negli header di Axios
  function setAuthToken(token) {
    if (token) {
      apiClient.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      localStorage.setItem('token', token);
    } else {
      delete apiClient.defaults.headers.common['Authorization'];
      localStorage.removeItem('token');
    }
  }
  
  // Recupera il token salvato all'avvio
  const savedToken = localStorage.getItem('token');
  if (savedToken) {
    setAuthToken(savedToken);
  }
  
  export default {
    getConversations(username) {
        return apiClient.get('/conversations', {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
          params: {
            username: username  // Inviando solo la stringa, senza l'oggetto
          }
        })
        .then(response => {
          console.log('Conversations fetched:', response.data);
          return response.data;
        })
        .catch(error => {
          console.error('Error fetching conversations:', error);
          if (error.response) {
            console.error('Response error:', error.response.data);
          } else {
            console.error('Request error:', error.message);
          }
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
    return apiClient.post(`/conversations/messages`, 
      { convoID, ...message }, 
      { headers: { Authorization: `Bearer ${token}` } }
    );
  },

  async login(credentials) {
    const response = await apiClient.post('/login', credentials);
    const token = response.data.token;

    if (response.data.user && typeof response.data.user.username === 'string') {
      try {
        response.data.user.username = JSON.parse(response.data.user.username).username;
      } catch (e) {
        console.error("Errore nel parsing del username:", e);
      }
    }

    setAuthToken(token);
    return response;
  },

  logout() {
    setAuthToken(null);
  }
};
