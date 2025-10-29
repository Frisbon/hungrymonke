<template>
  <div class="forwarding-list">
    <h3>Choose where to forward the message to:</h3>

    <div
      v-for="(c, index) in convertedConvos"
      :key="c.convoid || index"
      class="forward-tile"
      role="button"

      @click="chooseConvo(c)"
      @keydown.enter="chooseConvo(c)"
      @keydown.space.prevent="chooseConvo(c)"
    >
      <img
        class="pfp"
        :src="(c.chatPic && c.chatPicType) ? ('data:' + c.chatPicType + ';base64,' + c.chatPic) : 'https://i.imgur.com/D95gXlb.png'"
        @error="e => (e.target.src = 'https://i.imgur.com/D95gXlb.png')"
      />

      <h4 class="chatName">{{ c.chatName }}</h4>
    </div>
  </div>

  <div class="new-user-forward">
  <h4>... or forward to a new user</h4>
  <form @submit.prevent="createAndChoose" class="inline-form">
    <input
      class="user-search"
      type="text"
      v-model="newUsername"
      placeholder="Enter Username..."
      aria-label="Username"
    />
    <button type="submit" :disabled="creating || !newUsername"> Create and Forward</button>
  </form>
  <p v-if="createErr" class="err">{{ createErr }}</p>
</div>

</template>

<script>
import api from '@/api';


export default {
  name: 'ForwardingConvoList',
  props: { currentUser: String, convertedConvos: Array },
  data() {
    return {
      newUsername: '',
      creating: false,
      createErr: ''
    };
  },
  methods: {
    chooseConvo(convo) {
      this.$emit('forwardToConvo', convo);
    },
    async createAndChoose() {
      const u = (this.newUsername || '').trim();
      if (!u) return;
      this.creating = true; this.createErr = '';
      try {
        // crea (o recupera) la privata esistente
        const res = await api.startPrivateConvo(u);
        const id = res?.data?.convoID;
        if (!id) {
          this.createErr = (res?.data?.error || 'Cannot create conversation');
          return;
        }
        // costruisco un oggetto "convo" minimale compatibile con App.vue
        const newConvo = { convoid: id, chatName: u, chatPic: '', chatPicType: '' };
        this.$emit('forwardToConvo', newConvo);
      } catch (e) {
        this.createErr = e?.response?.data?.error || e?.message || 'Error creating convo';
      } finally {
        this.creating = false;
      }
    }
  }
}
</script>

<style scoped>
.forwarding-list{
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  max-height: 75vh;
  overflow-y: auto;
  padding-right: 6px; /* per evitare scrollbar sopra i tile */
}

.pfp{
  width: 60px;
  height: 60px;
  border-radius: 50%;
  overflow: hidden;
  display: inline-block;
  border: 1px solid #ccc;
  object-fit: cover;
}

.chatName{
  font-weight: 600;
  margin: 0; 
}


.forward-tile{
  display: flex;
  align-items: center;
  gap: 12px;

  width: 100%;
  max-width: 360px;          /* simile alla colonna sinistra */
  padding: 10px 14px;
  border: 1px solid #e2e8f0; /* grigio chiaro */
  border-radius: 12px;
  background: #fff;
  cursor: pointer;

  transition: background-color .15s ease, box-shadow .15s ease, border-color .15s ease, transform .02s ease;
}

.forward-tile:hover,
.forward-tile:focus{
  background: #f3f4f6;       /* bg grigio */
  border-color: #cbd5e1;
  box-shadow: 0 0 0 3px #e5e7eb inset; /* outline morbido */
  outline: none;
}

.forward-tile:active{
  transform: translateY(1px);
}

.new-user-forward { margin-top: 12px; width: 100%; max-width: 360px; align-self: center;}
.inline-form { display: flex; gap: 8px; }
.user-search{
  flex: 1 1 auto; padding: 8px 12px; border: 1px solid #e2e8f0; border-radius: 10px;
}
.user-search:focus{ outline: none; box-shadow: 0 0 0 3px #e5e7eb inset; border-color: #cbd5e1; }
.err{ color:#b91c1c; margin-top:6px; }



</style>
