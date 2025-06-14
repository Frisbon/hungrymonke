<template>
    

    <h3>Choose where to forward the message to:</h3>
    <div class="convoContainer"
        v-for="(c, index) in this.convertedConvos"
        :key="c.convoid || index"
        @click="chooseConvo(c)">

              <img class="pfp" v-if="c.chatPic != null && c.chatPic != 'https:// i.imgur.com/D95gXlb.png' && c.chatPic != ''"
              :src="'data:' + c.chatPicType + ';base64,' + c.chatPic">

              <img class="pfp" v-else :src="'https:// i.imgur.com/D95gXlb.png'">
              <h3 class='chatName'>{{ c.chatName }}</h3>

    </div>

</template>

<script>

 export default {
   name: 'ForwardingConvoList', //  Meglio usare lo stesso nome usato nel genitore per chiarezza

   props: {
     currentUser: String,
     convertedConvos: Array
   },

   data() {
     return {
     };
   },

   methods: {
       chooseConvo(convo) { //  Ricevi l'oggetto conversazione intero
           console.log("You've chosen the convo: ", convo);
           this.$emit('forwardToConvo', convo);
       },

   },

   //  mounted DEVE essere una funzione
   mounted() {
      console.log('ForwardingConvoList mounted.');

   }
}
</script>

<style>


.pfp {
  /* Make the element a square to ensure a perfect circle */
  width: 60px; /* Adjust the size as needed */
  height: 60px; /* Should be the same as the width */
  border-radius: 50%; /* This makes the element circular */
  overflow: hidden; /* Clips content that goes outside the circle */
  display: inline-block; /* Allows multiple profile pictures to sit on the same line */
  margin: 10px 0px;
  border: 1px solid #ccc;
  object-fit: cover

}

.convoContainer{
    align-self: center;
    width: 40%;
    display: flex;
    border: 1px solid #ccc;
    margin: 0px 10px 10px;
    border-radius: 10px;
    padding: 0 15px;
    justify-content: space-between;
    align-items: center;
}

</style>