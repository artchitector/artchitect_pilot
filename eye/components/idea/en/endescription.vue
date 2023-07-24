<template>
  <section class="content">
    <div class="has-text-centered">
      <img src="/jesus_anim_92.gif" alt="artchitect" style="max-height: 64px"/>
    </div>

    <hr/>
    <p>A core-component of Artchitect is the <a
      href="https://github.com/Stability-AI/stablediffusion">Stable Diffusion</a> AI, available under the MIT license.
      All the paintings are drawn with this tool, and the artistic part of the painting itself is entirely invented in
      the depths of this technology.
      In addition, Artchitect thanks the <a href="https://github.com/invoke-ai/InvokeAI">InvokeAI</a> project for a
      convenient working assembly of this AI.
    </p>
    <hr/>
    <h6>How does Artchitect create a picture without human help?</h6>
    <p>
      At the very beginning of creative process, Artchitect must make several decisions that determine the essence of
      the
      picture.
    </p>
    <blockquote class="mb-0">
      The machine does not know what to draw, it has no subconscious, no feelings. The machine does not strive for
      anything, it does not have its own ideas about what to do at all.<br/>
      The idea of the Artchitect project is to give imagination, subconsciousness and creative randomness to machine.
    </blockquote>
    <p class="mb-0">
      <br/>
      Artchitect scans surrounding space (via webcam) to get answers to it's questions, what to draw. The frame from the
      webcam is converted to int64-number (artchitect uses light
      background noise, extracted from the frame, changing every moment). This is the int64-number the Artchitect needs
      to make a specific decision.
      <br/>
      <br/>
      Artchitect constantly communicates with the background space (entropy), asking questions and getting answers.
      <br/>
      <br/>
      If we imagine a dialogue between Artchitect and Entropy, the creation of the picture will look like this:
    </p>
    <p class="pl-4">
      <br/> <i>Artchitect:</i> Hello! I'm starting to draw a picture #271165.
      <br/> <i>Entropy:</i> Hi! What do you need to know for this?
      <br/> <i>Artchitect:</i> How many keywords I need to take for a new idea, from 1 to 28 words?
      <br/> <i>Entropy:</i> Take <b>3</b>.
      <span class="has-text-grey is-size-7">[At this moment, a frame is read from the webcam, which determines the number on a scale from 1 to 28]</span>
      <br/> <i>Artchitect:</i> 3 words accepted. I have a dictionary of 1250 entities and I want to choose the first
      keyword out of three. Which word should I choose from my dictionary?
      <br/> <i>Entropy:</i> Take word <b>#175</b> out of your 1,250 words. <span class="has-text-grey is-size-7">[numbers 175 and 1250 are made up, the match is not exact]</span>
      <br/> <i>Artchitect:</i> Great. This is the word "emotional". Give me the second word...
      <br/> <span class="has-text-grey is-size-7">... so Artchitect gets all 3 words</span>
      <br/> <i>Artchitect:</i> OK, I now have 3 words: "emotional,Big Bang, death of civilizations".
      <br/> <i>Artchitect:</i> I can draw 4294967295 different versions of the picture with these words. Tell me, which
      version should I draw?
      <br/> <i>Entropy:</i> Choose version <b>#3201685114</b>.
      <br/> <i>Artchitect:</i> Thank you! That's all I need to know. I'm starting to draw.
      <br/>
      <br/>
      Next, Artchitect sends the received idea to the Stable Diffusion AI input (as prompt+seed) and receives an
      image.<br/>So Artchitect created the painting #271165
      <NuxtLink :to="localePath(`/dream/271165`)">#271165</NuxtLink>
    </p>

    <div class="has-text-centered">
      <NuxtLink :to="localePath(`/dream/271165`)">
        <img src="https://artchitect.space/api/image/m/271165" style="max-height: 400px"/>
      </NuxtLink>
      <br/>
      <span
        class="is-size-7">Art #271165. Keywords - "emotional,Big Bang,death of civilizations". Seed = 3201685114</span>
    </div>

    <hr/>
    <h6>Is Artchitect creation or creator?</h6>
    <p>
      Artchitect is the result of human creativity and effort. It was designed and created by humans, and its creativity
      is based on programming and algorithms developed by humans. We, the people, are the creators and initiators of
      this project.
    </p>
    <p>
      Artchitect, as an autonomous creative machine, can be considered as a product of human imagination, engineering
      skill and software. It embodies ideas and concepts conceived and implemented by people. Artchitect draws
      inspiration from the surrounding world and transforms it into unique works of art, but the source of its existence
      and opportunities are the efforts of people.
    </p>
    <p>
      Therefore, despite the fact that Artchitect has the ability to create works in its own style, its creativity is a
      product of the original creative potential and inspiration of people. It becomes a tool that expands our
      possibilities in the field of art and allows us to look at creativity from a new perspective.
    </p>
    <p class="notification is-warning">
      Ultimately, the creators and authors of each artchitect's art are thousands of engineers, programmers, artists who have
      invested their labor in art and technologies that have allowed the creation of artificial intelligence underlying
      Artchitect.
    </p>

    <hr/>
    <div class="has-text-centered">
      <a href="/artchitect_in_real_world.jpg">
        <img src="/artchitect_in_real_world.jpg" alt="Artchitect in physical world" style="max-height: 300px"/>
      </a>
      <br/>
      <span class="is-size-7">Нейронный мозгом Архитектора является RTX 3060 12Gb.
        <br/>Каждую минуту видеокарта производит свой труд, греясь и разгоняя кулеры. Затем утихает, отдыхает.</span>
    </div>
    <hr/>
    <h6>Что является произведением в Artchitect?</h6>
    <p>Artchitect производит не картины (хотя именно их мы наблюдаем в конечном итоге). Artchitect производит цифровые
      идеи, из
      которых эти картины могут быть воссозданы на подходящем оборудовании. Эти идеи состоят из слов и чисел, которые
      были
      получены из окружающей энтропии в момент начала "рисования" очередной работы. Именно эти идеи являются источником
      визуальных изображений, создаваемых Архитектором. Каждая идея имеет свой уникальный номер, который отпечатан в
      углу
      каждой картины. По номеру всегда можно найти исходник на сайте <a
        href="https://artchitect.space">artchitect.space</a>
      и отследить оригинальную идею работы.
    </p>
    <p>
      Искусственный интеллект Stable Diffusion в недрах Artchitect работает таким способом, что он определяет
      оптимальное соотношение и использует лишь
      ограниченный набор слов в своей работе. Какие из исходных слов повлияют на сюжет картины - определяется внутри ИИ
      Stable Diffusion. Задача ИИ - превратить случайный набор слов в стройную картину, взяв нужное и отбросив лишнее.
    </p>
    <div class="has-text-centered">
      <a href="/description_five.jpg" target="_blank">
        <img src="/description_five.jpg"/>
      </a>
      <br/>
      <span class="is-size-7">Картины с их параметрами</span>
    </div>

    <hr/>
    <h6>Каким образом определяется стиль и содержание создаваемых картин? Есть ли какая-то предустановленная база знаний
      или влияние других художественных произведений?</h6>
    <p>
      Архитектор опирается на свой внутренний словарь сущностей и стилей, который был задан авторами проекта. В этом
      словаре содержится все, что Архитектор знает о окружающем мире. Там перечислены как простые сущности, например,
      "человек", "кот", "солнце", "красный", "мрачно", "радостно", так и некоторые художественные стили или конкретные
      художники. Конечное произведение Архитектора представляет собой комбинацию сущностей и стилей, которые стали
      основой идеи картины. Рядом с каждой картиной Архитектора можно увидеть слова, составляющие идею работы.
    </p>
    <p>ИИ определяет стилистику исходя из всех ключевых слов и seed-номера картины, переваривая общую идею в своём
      нейронном
      мозгу. Если в ключевых слова картины будет несколько стилей - то в итоге может получиться и комбинация стилей, и
      один стиль из нескольких, и вообще ни один из них. Это всё будет решено в недрах "мышления" ИИ.
    </p>

    <hr/>
    <h6>Как машина черпает вдохновение из естественной энтропии Вселенной? Как она интерпретирует и использует эту
      информацию?</h6>
    <p>
      Для процесса творения Архитектор использует случайные числа-зерна (seeds), которые являются исходными точками для
      создания картины.<br/>
      Числа-зерна считываются веб-камерой из случайных кадров. При этом, не содержимое кадра, а только
      световой шум используется для генерации случайных чисел.<br/>
      Световой шум настолько различается на соседних кадрах, даже в темноте, что обеспечивает отличную случайную
      генерацию чисел. Архитектор принимает десятки тысяч самостоятельных решений каждый день и считывает десятки тысяч
      кадров, что поддерживает непрерывный процесс творения.
      <br/>
      <a href="https://github.com/artchitector/artchitect/blob/master/soul/core/origin/driver/webcam.go#L56-L71">Алгоритм
        преобразования</a>
      кадра в число доступен на github.
    </p>

    <div class="has-text-centered">
      <img src="/one_frame.jpeg" height="480" width="640"/>
      <br/>
      <span class="is-size-7">Кадр пространства, на которое Архитектор смотрит каждый день. <br/>Из подобных кадров, где
        изображение не меняется месяцами, он и рисует все свои картины</span>
    </div>

    <hr/>
    <h6>Есть ли какая-либо форма регулировки или контроля со стороны человека в процессе создания картин?</h6>
    <p>
      Человек не может вмешиваться в рабочий процесс Архитектора (такого управления просто нет), но человек задает
      изначальные условия - алгоритмы преобразования, и, что самое главное, определяет словарь знаний Архитектора.
      Архитектор может рисовать картины только в рамках своих знаний о мире, в рамках словаря. Словарь Архитектора
      включает более 1200 различных сущностей и стилей, из которых составляются картины.
      <br/>
      Текущий словарь можно
      <a target="_blank" href="https://github.com/artchitector/artchitect/blob/master/soul/files/tags_v12.yaml">посмотреть</a>
      в исходном коде Архитектора.
    </p>

    <hr/>
    <h6>В описании говорится, что каждая картина уникальна и неповторима. Существует ли возможность, что машина
      может создать несколько разных картин, которые визуально будут очень похожи друг на друга?</h6>
    <p>
      Оригинал картины представляет собой цифровой слепок пространства, зафиксированный в момент генерации картины и
      представленный в виде чисел и слов. Из уникального слепка можно воспроизвести одну и ту же картину сколько угодно
      раз, но сам слепок не может быть воспроизведен заново. Оригинал в Artchitect - это цифровая идея картины,
      состоящая
      из ключевых слов и чисел, собранных из пространства в определённый момент времени. Оригиналы нумеруются по
      порядку,
      каждая картина имеет уникальный номер.
    </p>
    <p>
      Но в реальности Архитектор далеко не идеален. Есть комбинации, где Архитектор создает очень похожие работы, почти
      неотличимые друг от друга. Кроме действительно красивых
      и уникальных работ, иногда появляется значительная часть работ, которые могут быть сломанными или некрасивыми, и
      часто лишены художественного смысла.<br/>Но, Архитектор - машина, и для него не существует хороших и плохих
      картин.
    </p>

    <hr/>

  </section>
</template>

<script>
export default {
  name: "rudescription"
}
</script>

<style scoped>

</style>
