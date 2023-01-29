<template>
  <div class="content">
    <h1 class="is-size-4">Request to the Universe / Digital prayer</h1>
    <p>
      Here you can send request to the Universe and receive answer in the form of your personal card.
    </p>
    <blockquote>
      Then he [Jesus] touched their eyes and said, “According to your faith let it be done to you";<br/>
      <i><a href="https://www.biblegateway.com/verse/en/Matthew%209%3A29" target="_blank">Matthew 9:29</a></i>
    </blockquote>
    <p>
      Main idea of Artchitect is that it interconnected with Quantize Fields, which controlled directly by Universe
      (more info in
      <NuxtLink to="/idea">idea section</NuxtLink>
      ).
      Universe/God absolutely knows that exactly You will send this request, and Universe/God can create next card
      specially for You, as answer for you. The idea of the card may not be clear, but you should think about it.
    </p>
    <p>Try it. Maybe the answer is closer than it seems.</p>
    <hr/>
    <p><b>First step</b>: think carefully about you request/pray, that you want to send to the Universe. Concentrate
      your mind on your request.</p>
    <p><b>Second step</b> (optional): describe your request in the text field. This will help you better understand
      and systematize the common idea of the request. Alternatively, you can write your request on a piece
      of paper or think it over in your head. The universe knows in advance and better than you what you are asking for.
      <br/>
      <b>You keep it brief!</b>
    </p>
    <p>
      <textarea class="textarea"
                rows="2"
                :disabled="loading"
                :placeholder="`Type your request message. Write carefully. Between you and Universe. Secure - data not being send anywhere.`"></textarea>
    </p>
    <p>
      <b>Third step</b>: finalize request with "Pray" button. This will erase text from textarea and start card creation
      process.
      You will see your result when it finished.
    </p>
    <div class="notification">
      <b>Important</b>: nobody will know, what you asked and what you received. Data about request not sends anywhere.
      But your personal card will appear in common list of cards at
      <NuxtLink to="/">artchitect.space</NuxtLink>
      .
      There will no any links to your person (no any personal data saved anywhere).
    </div>
    <p class="has-text-centered">
      <button class="button is-primary" :disabled="loading" @click.prevent="submit()">Submit your request (Amen!)
      </button>
    </p>
    <p v-if="error" class="notification is-danger">{{ error }}</p>
    <p v-if="loading" class="has-text-centered">
      <loader/>
    </p>
  </div>
</template>

<script>
export default {
  name: "prayer",
  data() {
    const randomPassword = (Math.random() + 1).toString(36).substring(7);
    return {
      password: randomPassword,
      loading: false,
      error: null,
    }
  },
  methods: {
    async submit() {
      try {
        this.loading = true
        localStorage.setItem("last_pray_password", this.password);  // only emitter of pray can see answer (temporary password needed)
        const result = await this.$axios.$post("/pray", {
          password: this.password
        })
        await this.$router.push(`/pray/${result}`)
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>

</style>
