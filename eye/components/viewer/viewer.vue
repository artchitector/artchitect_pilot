<template>
  <div v-if="isVisible" class="modal-container">
    <div class="background" @click="close()"></div>
    <div class="control-prev" v-if="hasPrev">
      <a href="#" @click.prevent="prev()"><</a>
    </div>
    <div class="control-next" v-if="hasNext">
      <a href="#" @click.prevent="next()">></a>
    </div>
    <div class="header">
      <h1 class="is-size-4" v-if="card">
        <NuxtLink :to="`/card/${card.ID}`">
          {{ $t('card') }} #{{ card.ID }}
        </NuxtLink>
      </h1>
      <p v-if="list.length > 1" class="has-text-centered">
        {{ index + 1 }} / {{ list.length }}
      </p>
    </div>
    <div class="img">
      <loader v-if="loading"/>
      <div v-else-if="error">
        <div class="notification is-danger">
          <p>Ошибка:</p>
          <p>{{ error }}</p>
          <p>Попробуйте позже, сейчас Артхитектору плохо</p>
        </div>
        <div class="has-text-centered">
          <button class="button" @click="close()">Закрыть</button>
        </div>
      </div>
      <!-- Main image here-->
      <img v-else-if="card" :src="`/api/image/f/${card.ID}`"/>
      <!--      -->
    </div>
    <div class="tags">
      <template v-if="card">
        <p>{{ $t('created') }}: {{ formatDate(card.CreatedAt) }}</p>
        <p>{{ $t('seed') }}: {{ card.Spell.Seed }}</p>
        <p class="is-size-7">{{ $t('tags') }}: {{ card.Spell.Tags }}</p>
      </template>
    </div>
  </div>
</template>
<script>
import Loader from "@/components/loader";
import moment from "moment/moment";

export default {
  components: {Loader},
  data () {
    return {
      isVisible: false,
      loading: false,
      list: [], // all cards
      card_id: null, // current card_id
      index: null, // current card index in list
      card: null, // current loaded card
      error: null,
    }
  },
  computed: {
    hasPrev () {
      return this.list.length > 1 && this.index > 0
    },
    hasNext () {
      return this.list.length > 1 && this.index < this.list.length - 1
    }
  },
  methods: {
    show (list, card_id) {
      this.isVisible = true
      this.list = list
      this.card_id = card_id
      this.index = this.list.indexOf(card_id)
      this.load()
      window.addEventListener('keyup', this.onGlobalKey)
    },
    async load () {
      if (!this.card_id) {
        return
      }
      this.loading = true
      try {
        this.card = await this.$axios.$get(`/card/${this.card_id}`)
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    formatDate (date) {
      // TODO need make global date helper and use it everywhere
      return moment(date).format("YYYY MMM Do HH:mm:ss")
    },
    close () {
      this.isVisible = false
      this.list = null
      this.card_id = null
      this.card = null
      window.removeEventListener('keyup', this.onGlobalKey)
    },
    onGlobalKey (e) {
      if (e.key === 'Escape') {
        this.close()
      } else if (e.key === 'ArrowLeft') {
        this.prev()
      } else if (e.key === 'ArrowRight') {
        this.next()
      }
    },
    setIndex (index) {
      this.index = index
      this.card_id = this.list[index]
      this.card = null
      this.load()
    },
    prev () {
      if (!this.hasPrev) {
        return
      }
      this.setIndex(this.index - 1)
    },
    next () {
      if (!this.hasNext) {
        return
      }
      this.setIndex(this.index + 1)
    },
  }
}
</script>
<style lang="scss">
.modal-container {
  padding: 20px;
  position: fixed;
  z-index: 1;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #d4d1c3;
  gap: 10px;

  .control-prev {
    position: fixed;
    left: calc(10vw);
    top: 50%;
    font-size: 50px;
    z-index: 3;
    font-weight: bolder;
  }
  .control-next {
    position: fixed;
    right: calc(10vw);;
    top: 50%;
    font-size: 50px;
    z-index: 3;
    font-weight: bolder;
  }

  .background {
    bottom: 0;
    left: 0;
    position: absolute;
    right: 0;
    top: 0;
    background-color: rgba(0, 0, 0, 0.8);
  }

  .header {
    z-index: 2;

    h1 {
      background-color: rgba(0, 0, 0, 0.5);
    }
  }

  .img {
    z-index: 2;
    max-height: 100%;
    overflow: hidden;

    img {
      max-height: 100%;
    }
  }

  .tags {
    z-index: 2;
    display: block;
    max-width: calc(60vw);
    background-color: rgba(0, 0, 0, 0.5);
  }
}
</style>
