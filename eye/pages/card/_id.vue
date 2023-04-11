<i18n>
{
  "en": {
    "title": "Artchitect - dream",
    "card": "Dream",
    "loading": "loading...",
    "error": "error: ",
    "version": "version",
    "seed": "Seed",
    "tags": "entities",
    "full_size": "view full size"
  },
  "ru": {
    "title": "Artchitect - сон",
    "card": "Сон",
    "loading": "загрузка...",
    "error": "ошибка: ",
    "version": "версия",
    "seed": "Зерно",
    "tags": "сущности",
    "full_size": "смотреть в полном размере"}
}
</i18n>
<template>
  <section>
    <div class="notification is-primary" v-if="$fetchState.pending">
      {{$t('loading')}}
    </div>
    <div class="notification is-danger" v-if="$fetchState.error">
      {{$t('error')}} {{ $fetchState.error.message }}
    </div>
    <div v-else-if="card">
      <h1 class="is-size-2 has-text-centered">{{ $t('card') }} #{{ card.ID }}</h1>
      <p><span class="tag is-primary is-light">{{$t('version')}} {{ card.Version }}</span></p>
      <p>{{ created }}</p>
      <p>{{$t('seed')}} = {{ card.Spell.Seed }}</p>
      <p class="tags">{{$t('tags')}} = <i>{{ card.Spell.Tags }}</i></p>
      <p class="has-text-centered">
        <a :href="fullSizeUrl" target="_blank" class="is-size-7">{{$t('full_size')}}</a>
      </p>
      <img :src="`/api/image/f/${card.ID}`"/>
    </div>
  </section>

</template>
<script>
import moment from "moment"

export default {
  head() {
    return {
      title: this.$t('title')+ ` #${this.$route.params.id}`
    }
  },
  data() {
    return {
      card: null
    }
  },
  computed: {
    created() {
      return moment(this.card.CreatedAt).format("YYYY MMM Do HH:mm:ss")
    },
    fullSizeUrl() {
      return `${process.env.STORAGE_URL}/cards/card-${this.card.ID}.jpg`
    }
  },
  async fetch() {
    const id = parseInt(this.$route.params.id);
    if (!id) {
      throw "id must be positive integer"
    }
    this.card = await this.$axios.$get(`/card/${id}`)
  }
}
</script>
<style lang="scss" scoped>
p.tags {
  word-wrap: break-word;
  word-break: break-all;
  overflow: hidden;
}
</style>
