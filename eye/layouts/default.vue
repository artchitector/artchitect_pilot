<template>
  <div>
    <template v-html="yandexMetrica"></template>
    <div class="has-text-centered is-header pt-1 pb-2">
      <div>
        <NuxtLink :to="{name: 'index'}">
          <h1 class="is-size-3 has-text-weight-bold">
            artchitect.space
          </h1>
        </NuxtLink>
      </div>
      <div>
        <NuxtLink :to="{name: 'index'}">main</NuxtLink>
        <span> - </span>
        <NuxtLink :to="{name: 'lottery'}">lottery</NuxtLink>
        <span> - </span>
        <NuxtLink :to="{name: 'selection'}">selection</NuxtLink>
      </div>
    </div>

    <div class="container is-max-desktop">
      <Nuxt/>
    </div>

    <footer class="footer">
      <p class="has-text-centered">
        github: <a href="https://github.com/artchitector/artchitect" target="_blank">https://github.com/artchitector/artchitect</a><br/>
        telegram: <a href="https://t.me/artchitect_space" target="_blank">t.me/artchitect_space</a><br/>
        twitter: <a href="https://twitter.com/artchitector" target="_blank">@artchitector</a>
      </p>
    </footer>
  </div>
</template>
<script>
export default {
  computed: {
    yandexMetrica() {
      return `
      <!-- Yandex.Metrika counter -->
    <script type="text/javascript" >
      (function(m,e,t,r,i,k,a){m[i]=m[i]||function(){(m[i].a=m[i].a||[]).push(arguments)};
        m[i].l=1*new Date();
        for (var j = 0; j < document.scripts.length; j++) {if (document.scripts[j].src === r) { return; }}
        k=e.createElement(t),a=e.getElementsByTagName(t)[0],k.async=1,k.src=r,a.parentNode.insertBefore(k,a)})
      (window, document, "script", "https://mc.yandex.ru/metrika/tag.js", "ym");

      ym(92022377, "init", {
        clickmap:true,
        trackLinks:true,
        accurateTrackBounce:true
      });
    <\/script>
<noscript><div><img src="https://mc.yandex.ru/watch/92022377" style="position:absolute; left:-9999px;" alt="" /></div></noscript>
<!-- /Yandex.Metrika counter -->
      `
    }
  }
}
</script>
<style lang="scss">
.is-header {
  background-color: #2d2d2d;

  h1 {
    color: #d4d1c3;
  }
}

footer.footer {
  background-color: #2d2d2d;
  color: hsl(0, 0%, 96%);
}
</style>
