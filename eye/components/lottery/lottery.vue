<template>
  <div class="box">
    <h3 class="is-size-4">Lottery "{{ lottery.Name }}" - {{ lottery.State }}</h3>
    <p>Start time: <i>{{ formatDate(lottery.StartTime) }}</i></p>
    <p>Collection period: <i>{{ formatDate(lottery.CollectPeriodStart) }} - {{
        formatDate(lottery.CollectPeriodEnd)
      }}</i></p>

    <div v-if="showWinners" class="has-text-centered box has-background-link-light">
      Winners ({{ lottery.Winners.length }} of total {{ lottery.TotalWinners }})<br/>
      <template v-for="winnerID in winners">
      <img v-if="winnerID === 0" class="mini-preview ml-1 mr-1" src="/in-progress-lottery.jpg"/>
      <a v-else :href="`/card/${winnerID}`" @click.prevent="select(winnerID)" class="winner-link">
        <img class="mini-preview ml-1 mr-1" :src="`/api/image/xs/${winnerID}`"/>
      </a>
      </template>
    </div>
    <viewer ref="viewer"/>
  </div>
</template>
<script>
import moment from "moment";

export default {
  props: ['lottery'],
  computed: {
    showWinners() {
      return this.lottery.State === "running" || this.lottery.State === "finished";
    },
    winners() {
      if (!this.lottery.TotalWinners || !this.lottery.Winners) {
        return [];
      }
      const winners = [];
      for (let i = 0; i < this.lottery.TotalWinners; i++) {
        if (i >= this.lottery.Winners.length) {
          winners.push(0);
        } else {
          winners.push(this.lottery.Winners[i]);
        }
      }
      return winners;
    }
  },
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
.winner-link {
  display: inline-block;
}
</style>
