<template>
  <div>
    <viewer ref="viewer" @liked="onLiked"/>
    <div class="columns" v-for="line in lines">
      <div class="column" v-for="card in line">
        <card-complex v-if="isComplex" :card="card" @select="select(card.ID)"/>
        <card-simple v-else :card-id="card" @select="select(card)"/>
      </div>
    </div>
    <div v-if="showLoadMore" class="has-text-centered">
      <button class="button" @click.prevent="showMore">show more...</button>
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
    'visibleCount', // how many cards show at start of component
  ],
  data() {
    return {
      currentVisible: -1,
    }
  },
  computed: {
    lines() {
      let cards = []
      if (this.currentVisible === -1) {
        cards = []
      } else if (this.currentVisible === 0) {
        cards = this.cards
      } else {
        cards = this.cards.slice(0, this.currentVisible)
      }
      const chunkSize = parseInt(this.cardsInColumn);
      const chunks = [];
      for (let i = 0; i < cards.length; i += chunkSize) {
        chunks.push(cards.slice(i, i + chunkSize));
      }
      return chunks;
    },
    isComplex() {
      return typeof this.cards[0] === 'object';
    },
    showLoadMore() {
      return this.currentVisible > 0 && this.currentVisible < this.cards.length;
    }
  },
  mounted() {
    if (!!this.visibleCount) {
      this.currentVisible = parseInt(this.visibleCount)
    } else {
      this.visibleCount = 0
    }
  },
  methods: {
    select(cardId) {
      const ids = [];
      const isComplex = this.isComplex
      this.cards.forEach((card) => {
        if (isComplex) {
          ids.push(card.ID)
        } else {
          ids.push(card)
        }
      })
      this.$refs.viewer.show(ids, cardId)
    },
    showMore() {
      this.currentVisible += parseInt(this.visibleCount)
    },
    onLiked(data) {
      this.$emit("liked", data)
    }
  }
}
</script>

<style scoped>

</style>
