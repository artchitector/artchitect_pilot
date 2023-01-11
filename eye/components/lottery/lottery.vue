<template>
  <div class="box">
    <h3 class="is-size-4">Lottery "{{ lottery.Name }}" - {{ lottery.State }}</h3>
    <p>Start time: {{ lottery.StartTime }}</p>
    <p>Collection period: {{ lottery.CollectPeriodStart }} - {{ lottery.CollectPeriodEnd }}</p>
    <div v-if="lottery.Winners && lottery.Winners.length" class="has-text-centered box has-background-link-light">
      Winners<br/>
      <a :href="`/api/painting/${cardID}`" v-for="cardID in lottery.Winners" target="_blank">
        <img class="mini-preview ml-1 mr-1" :src="`/api/painting/${cardID}`"/>
      </a>
    </div>
    <div v-if="isRunning && ours.length" v-for="tour in tours" class="has-text-centered box">
      tour <b>{{ tour.Name }}</b> winners (id={{ tour.ID }}) <br/>
      <a :href="`/api/painting/${cardID}`" v-for="cardID in tour.Winners" target="_blank">
        <img class="micro-preview ml-1 mr-1" :src="`/api/painting/${cardID}`"/>
      </a>
    </div>
  </div>
</template>
<script>
export default {
  props: ['lottery'],
  computed: {
    isRunning() {
      return this.lottery.State == 'running'
    },
    tours () {
      if (!this.lottery || !this.lottery.Tours || !this.lottery.Tours.length) {
        return []
      }
      return this.lottery.Tours.slice()
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
