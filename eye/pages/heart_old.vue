<i18n>
{
  "en": {
    "title": "Artchitect's heart ❤️"
  },
  "ru": {
    "title": "Сердце  Artchitect ❤️"
  }
}
</i18n>
<template>
  <div id="heart-main-backing">
<!--    <div class="heart-caption">-->
<!--      Hy heart is gonna blow from this Love!-->
<!--    </div>-->
    <header class="has-text-centered">
      <NuxtLink :to="localePath('index')">
      <img src="/icon64.png" alt="artchitect"/>
      </NuxtLink>
      <h1 class="is-size-3 has-text-weight-bold">
        <NuxtLink :to="localePath('index')">
        artchitect.space
        </NuxtLink>
      </h1>
    </header>
    <div id="main">
      <article>
        <heart/>
      </article>
<!--      <nav>Nav</nav>-->
<!--      <aside>Aside</aside>-->
    </div>
<!--    <footer>Footer</footer>-->
  </div>
</template>

<script>
import Heart from "@/components/big_heart/heart.vue"
export default {
  components: {Heart},
  layout: "empty",
  head() {
    return {
      title: this.$t('title')
    }
  }
}
</script>

<style lang="scss" scoped>
#heart-main-backing {
  background-color: #2d2d2d;
  display: flex;
  min-height: 100vh;
  flex-direction: column;
  margin: 0;
  color: #d4d1c3;
}
#main {
  display: flex;
  flex: 1;
}
header {
  padding: 1em 1em 0.3em 1em;
  background-color: #d4d1c3;
  //background-image: url("/heart/back1.jpg");
  background-size:100% 100%;
  color: #2d2d2d;
  img {
    margin: 0;
    box-shadow: 0 0 10px #2d2d2d;
  }
  h1 {
    margin: 0;
    margin-top: -35px;
    padding: 0;
    text-shadow: 0 0 3px #d4d1c3;
    a {
      color: #2d2d2d;
    }
  }
}
#main > article {
  flex: 1;
  order: 1;
}
#main > nav,
#main > aside {
  flex: 0 0 20vw;
}
#main > nav {
  background: #D7E8D4;
  order: 3;
}
#main > aside {
  background: beige;
  order: 2;
}
</style>
