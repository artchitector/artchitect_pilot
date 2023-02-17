<template>
  <div>
    <div class="columns" v-for="line in lines">
      <div class="column" v-for="hundred in line">
        <hundred :hundred="hundred"/>
      </div>
    </div>
  </div>
</template>

<script>
import Hundred from "@/components/search/hundred.vue";

export default {
  name: "hundreds-list",
  components: {Hundred},
  props: ["hundreds", "cardsInColumn", "visibleCount"],
  data() {
    return {
      currentVisible: -1
    }
  },
  computed: {
    lines() {
      let hundreds = []
      if (this.currentVisible === -1) {
        hundreds = []
      } else if (this.currentVisible === 0) {
        hundreds = this.hundreds
      } else {
        hundreds = this.hundreds.slice(0, this.currentVisible)
      }
      const chunkSize = parseInt(this.cardsInColumn)
      const chunks = [];
      for (let i = 0; i < hundreds.length; i += chunkSize) {
        console.log(`i: ${i}, hl: ${hundreds.length} cs: ${chunkSize}`)
        chunks.push(hundreds.slice(i, i + chunkSize))
      }
      console.log(chunks)
      return chunks
    }
  },
  mounted() {
    if (!!this.visibleCount) {
      this.currentVisible = parseInt(this.visibleCount)
    } else {
      this.currentVisible = 0
    }
  }
}
</script>

<style scoped>

</style>
