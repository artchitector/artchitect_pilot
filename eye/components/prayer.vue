<template>
  <div class="box">
    <div v-if="!visible" class="has-text-centered">
      <button class="button is-success" @click="visible=true">To pray!</button>
    </div>
    <div v-else>
      <textarea class="textarea" v-model="pray"
                placeholder="Place for pray message. Write carefully. Between you and God. Secure, data burns and not being send anywhere."></textarea>
      <div class="has-text-centered mt-2"><button class="button" @click="clear()">Pray! (burn text)</button></div>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      visible: false,
      pray: ""
    }
  },
  methods: {
    clear() {
      this.pray = ""
    }
  }
}
</script>
