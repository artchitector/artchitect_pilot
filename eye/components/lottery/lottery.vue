<template>
  <div class="box">
    <h3 class="is-size-4">Lottery "{{ lottery.Name }}" - {{ lottery.State }}</h3>
    <p>Start time: {{ lottery.StartTime }}</p>
    <p>Collection period: {{ lottery.CollectPeriodStart }} - {{ lottery.CollectPeriodEnd }}</p>
    <div v-if="lottery.Winners && lottery.Winners.length" class="has-text-centered box has-background-link-light">
      Winners ({{lottery.Winners.length}} of total {{lottery.TotalWinners}})<br/>
      <NuxtLink :to="`/card/${cardID}`" v-for="cardID in lottery.Winners" target="_blank">
        <img class="mini-preview ml-1 mr-1" :src="`/api/image/xs/${cardID}`"/>
      </NuxtLink>
    </div>
  </div>
</template>
<script>
export default {
  props: ['lottery']
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
