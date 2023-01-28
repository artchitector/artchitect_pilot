<template>
  <div>
    <viewer ref="viewer"/>
    <div class="columns" v-for="line in lines">
      <div class="column" v-for="card in line">
        <card-complex v-if="isComplex" :card="card" @select="select(card.ID)"/>
        <card-simple v-else :cardId="card"/>
      </div>
    </div>
  </div>
</template>

<script>
import CardComplex from "@/components/list/card/card-complex.vue";
import CardSimple from "@/components/list/card/card-simple.vue";

export default {
  name: "cardlist",
  components: {CardSimple, CardComplex},
  props: [
    'cards',
    'cardsInColumn',
    'cardSize',
  ],
  computed: {
    lines() {
      const chunkSize = parseInt(this.cardsInColumn);
      const chunks = [];
      for (let i = 0; i < this.cards.length; i += chunkSize) {
        chunks.push(this.cards.slice(i, i + chunkSize));
      }
      return chunks;
    },
    isComplex() {
      return typeof this.cards[0] === 'object';
    }
  },
  methods: {
    select(cardId) {
      const ids = [];
      this.cards.forEach((card) => ids.push(card.ID))
      this.$refs.viewer.show(ids, cardId)
    }
  }
}
</script>

<style scoped>

</style>
