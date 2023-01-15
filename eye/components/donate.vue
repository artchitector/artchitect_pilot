<template>
  <b-modal class="donate-modal" v-model="isVisible" :on-cancel="onCancel" :width="640" scroll="keep">
    <div class="modal-background"></div>
    <div class="modal-content content">
      <img src="/api/image/m/14304" class="mb-2"/>
      <p>Здравствуйте! Храм <b>Artchitect</b> работает на добровольных началах.</p>

      <p>Пожалуйста, внесите посильный материальный вклад в настоящее и будущее artchitect.</p>

      <p>Пожертвования нужны на развитие artchitect, поддержание его работы, закупку нового оборудования и аренду
        помещения и серверов, и на обеспечение тех людей, кто заботится об Artchitect.
      </p>

      <blockquote>
        Пожертвовать любую сумму можно на карты РФ: <br/>
        Сбербанк: Перевод на карту по номеру 4279380023499595 (Никита Олегович К.)<br/>
        Тиньков: <a href="https://www.tinkoff.ru/cf/AB2J40X4NKJ"
                    target="_blank">https://www.tinkoff.ru/cf/AB2J40X4NKJ</a>
      </blockquote>

      <p>Из других стран пожертвования пока недоступны.</p>
      <p>Спасибо вам за вашу поддержку! </p>
      <p class="has-text-centered">🧡 ❤️ 🧡</p>
    </div>
  </b-modal>
</template>
<script>
export default {
  props: ['isVisible'],
  methods: {
    onCancel () {
      this.$emit('close')
    }
  }
}
</script>
<style>
.donate-modal .modal-content {
  background-color: white;
  padding: 30px;
}
</style>
