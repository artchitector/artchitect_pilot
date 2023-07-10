<template>
  <div class="wrapper">
    <template v-if="message && message.CardID">
      <result :dream_id="message.CardID" :current-enjoy-time="message.CurrentEnjoyTime"
              :total-enjoy-time="message.EnjoyTime" :rnd-four="rndFour"/>
    </template>

    <template v-else>
      <div class="wrapper-cell">
        <progress-view v-if="message" :message="message"/>
      </div>
      <div class="wrapper-cell">
        <lastdream v-if="lastDreamID" :last="lastDreamID" :four="rndFour"/>
      </div>
    </template>

  </div>
</template>

<script>
import Lastdream from "@/components/flexheart/creation/lastdream.vue";
import ProgressView from "@/components/flexheart/creation/progress-view.vue";
import rnd from "@/components/big_heart/layout/creation/rnd.vue";
import Result from "@/components/flexheart/creation/result.vue";

export default {
  name: "creation",
  computed: {
    rnd() {
      return rnd
    }
  },
  components: {Result, Lastdream, ProgressView},
  data() {
    return {
      message: null,
      lastDreamID: null,
      rndFour: [],
    }
  },
  methods: {
    onMessage(channelName, msg) {
      if (channelName === 'creation') {
        this.message = msg
        if (this.message.PreviousCardID && this.message.PreviousCardID !== this.lastDreamID) {
          this.lastDreamID = this.message.PreviousCardID
        } else if (!this.message.PreviousCardID) {
          this.lastDreamID = null
        }
      } else if (channelName === 'heart') {
        if (msg.Rnd.length > 0) {
          this.rndFour = msg.Rnd
        }
      }
    }
  }
}
</script>

<style scoped lang="scss">
.wrapper {
  width: 100%;
  height: 100%;
  padding: 10px;

  letter-spacing: 1px;
  text-transform: uppercase;

  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: stretch;

  .wrapper-cell {
    align-self: center;
  }
}
</style>
