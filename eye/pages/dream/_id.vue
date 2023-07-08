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
    "full_size": "смотреть в полном размере"
  }
}
</i18n>
<template>
  <section>
    <div class="notification is-primary" v-if="$fetchState.pending">
      {{ $t('loading') }}
    </div>
    <div class="notification is-danger" v-if="$fetchState.error">
      {{ $t('error') }} {{ $fetchState.error.message }}
    </div>
    <div v-else-if="card">
      <h1 class="is-size-2 has-text-centered">{{ $t('card') }} #{{ card.ID }}</h1>
      <p><span class="tag is-primary is-light">{{ $t('version') }} {{ card.Version }}</span></p>
      <p>{{ created }}</p>
      <p>{{ $t('seed') }} = {{ card.Spell.Seed }}</p>
      <p class="tags">{{ $t('tags') }} = <i>{{ card.Spell.Tags }}</i></p>
      <p class="has-text-centered">
        <a :href="fullSizeUrl" target="_blank" class="is-size-7">{{ $t('full_size') }}</a>
      </p>
      <div class="image-container">
        <img :src="`/api/image/f/${card.ID}`"/>
        <div class="control-like">
          <font-awesome-icon v-if="liked && liked.error"
                             icon="fa-solid fa-triangle-exclamation"
                             :title="liked.error.message"/>
          <a v-else href="#" @click.prevent="like()">
            <font-awesome-icon v-if="!liked || !liked.liked" icon="fa-solid fa-heart" class="has-color-base"/>
            <font-awesome-icon v-else icon="fa-solid fa-heart" class="has-text-danger"/>
          </a>
        </div>
      </div>
    </div>
  </section>

</template>
<script>
import moment from "moment"

export default {
  head() {
    return {
      title: this.$t('title') + ` #${this.$route.params.id}`
    }
  },
  data() {
    return {
      card: null,
      likes: 0,
      liked: {
        error: null,
        liked: false,
      }
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
  methods: {
    async like() {
      try {
        let like = await this.$axios.$post("/like", {
          card_id: this.card.ID,
        })
        this.$emit('liked', like)
        this.liked = {
          id: like.ID,
          liked: like.Liked,
        };
        if (like.Liked) {
          this.card.Likes += 1
        } else {
          this.card.Likes -= 1
        }
      } catch (e) {
        console.error(e)
        this.liked = {
          error: e
        };
      }

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

.image-container {
  position: relative;

  .control-like {
    position: absolute;
    left: 50%;
    bottom: 20%;
    z-index: 3;
    margin-left: -20px;
    font-size: 48px;
    opacity: 0.7;
    filter: drop-shadow(0px 0px 8px rgba(255, 0, 0, 0.6));
  }
}
</style>
