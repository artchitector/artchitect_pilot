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
        <NuxtLink :to="localePath(`/dream/${message.PreviousCardID}`)">
          <img :src="`/api/image/f/${message.PreviousCardID}`"/>
        </NuxtLink>
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
  props: ["message", "heartState"]
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

}
</style>
