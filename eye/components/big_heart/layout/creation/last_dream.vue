<template>
  <div class="last-dream-box" v-if="message && message.PreviousCardID">
    <div class="wrapper">
      <div class="left">
        <rnd v-if="heartState && heartState.Rnd.length > 0" :card-id="heartState.Rnd[0]"/>
        <rnd v-if="heartState && heartState.Rnd.length > 1" :card-id="heartState.Rnd[1]"/>
      </div>
      <div class="center">
        <div class="mb-1">last dream was
          <NuxtLink :to="localePath(`/dream/${message.PreviousCardID}`)" class="has-text-info">
            #{{ message.PreviousCardID }}
          </NuxtLink>
        </div>
        <div class="last-dream-wrapper">
          <NuxtLink :to="localePath(`/dream/${message.PreviousCardID}`)">
            <img :src="`/api/image/f/${message.PreviousCardID}`"/>
          </NuxtLink>
          <div class="control-like">
            <font-awesome-icon v-if="liked && liked.error"
                               icon="fa-solid fa-triangle-exclamation"
                               :title="liked.error.message"/>
            <a v-else href="#" @click.prevent="like()">
              <font-awesome-icon v-if="!liked || !liked.liked" icon="fa-solid fa-heart" class="has-color-base"/>
              <font-awesome-icon v-else icon="fa-solid fa-heart" class="has-text-danger"/>
            </a>
          </div>
        </div>
      </div>
      <div class="right">
        <rnd v-if="heartState && heartState.Rnd.length > 2" :card-id="heartState.Rnd[2]"/>
        <rnd v-if="heartState && heartState.Rnd.length > 3" :card-id="heartState.Rnd[3]"/>
      </div>
    </div>
  </div>
</template>

<script>
import Rnd from "@/components/big_heart/layout/creation/rnd.vue";

export default {
  components: {Rnd},
  props: ["message", "heartState"],
  data() {
    return {
      liked: {
        liked: false,
        error: null,
      }
    }
  },
  mounted() {
    this.initLiked()
  },
  methods: {
    async like() {
      try {
        let like = await this.$axios.$post("/like", {
          card_id: this.message.PreviousCardID,
        })
        this.$emit('liked', like)
        this.liked = {
          id: like.ID,
          liked: like.Liked,
        };
      } catch (e) {
        console.error(e)
        this.liked = {
          error: e
        };
      }

    },
    async initLiked() {
      const cardID = this.message.PreviousCardID
      try {
        let like = await this.$axios.$get(`/liked/${cardID}`)
        this.liked.liked = like.Liked
      } catch (e) {
        console.error(e)
        this.liked.error = e
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.last-dream-box {
  padding: 0 10px;
  background-color: rgba(0, 0, 0, 0.1);
  position: absolute;
  bottom: 0;
  left: 50%;
  width: 100%;
  margin-left: -50%;
  margin-bottom: 20px;

  .wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    max-width: 700px;
    margin: auto;

    .center {
      display: block;
      flex-basis: 0;
      flex-grow: 3;
      flex-shrink: 1;
      max-width: 400px;
      padding: 0 0.75rem 0.75rem;
    }

    .left, .right {
      display: block;
      flex-grow: 1;
      width: 20px;
      padding: 40px 0.4rem 0.4rem;
    }
  }

  .last-dream-wrapper {
    position: relative;

    .control-like {
      position: absolute;
      left: 50%;
      bottom: 20%;
      z-index: 3;
      margin-left: -20px;
      font-size: 48px;
      opacity: 0.7;
      filter: drop-shadow(0px 0px 8px rgba(255, 0, 0, 0.6));
    }
  }
}
</style>
