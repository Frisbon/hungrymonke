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
</template>

<script>
export default {
  name: 'ForwardingConvoList',
  props: {
    currentUser: String,
    convertedConvos: Array
  },
  methods: {
    chooseConvo(convo) {
      console.log("You've chosen the convo: ", convo);
      this.$emit('forwardToConvo', convo);
    },
    close() { this.$emit('close'); }
  },
  mounted() {
    console.log('ForwardingConvoList mounted.');
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
</style>
