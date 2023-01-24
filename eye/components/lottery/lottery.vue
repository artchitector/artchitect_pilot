<template>
  <div class="box">
    <h3 class="is-size-4">Lottery "{{ lottery.Name }}" - {{ lottery.State }}</h3>
    <p>Start time: <i>{{ formatDate(lottery.StartTime) }}</i></p>
    <p>Collection period: <i>{{ formatDate(lottery.CollectPeriodStart) }} - {{ formatDate(lottery.CollectPeriodEnd) }}</i></p>
    <div v-if="lottery.Winners && lottery.Winners.length" class="has-text-centered box has-background-link-light">
      Winners ({{lottery.Winners.length}} of total {{lottery.TotalWinners}})<br/>
<!--      <NuxtLink :to="`/card/${cardID}`" v-for="cardID in lottery.Winners" target="_blank">-->
<!--        <img class="mini-preview ml-1 mr-1" :src="`/api/image/xs/${cardID}`"/>-->
<!--      </NuxtLink>-->
      <a :href="`/card/${cardID}`" @click.prevent="select(cardID)" v-for="cardID in lottery.Winners">
        <img class="mini-preview ml-1 mr-1" :src="`/api/image/xs/${cardID}`"/>
      </a>
    </div>
    <viewer ref="viewer"/>
  </div>
</template>
<script>
import moment from "moment";
export default {
  props: ['lottery'],
  methods: {
    formatDate(date) {
      // TODO need make global date helper and use it everywhere
      return moment(date).format("YYYY MMM Do HH:mm:ss")
    },
    select(id) {
      this.$refs.viewer.show(this.lottery.Winners, id);
    }
  }
}
</script>
<style>
img.mini-preview {
  max-height: 100px;
}

img.micro-preview {
  max-height: 60px;
}
</style>
