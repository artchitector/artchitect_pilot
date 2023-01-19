<i18n>
{
  "en": {
    "how": "How is it technically arranged",
    "zero": "As the result, we expect to get millions of artistic paintings created according to God's Plan. We will try to see this idea, if there is one plan here. Or we won't see it if it's not there.",
    "one": "The picture itself is technically drawn by modern artificial intelligence, in which an infinite (this is important) number of creations can be created. The Artchitect uses Stable Diffusion AI. God can choose any variant of the picture He wants. He knows all the possible variants in advance.",
    "two": "To generate picture Stable Diffusion need specified set of text tags and a special number - the Seed. To make system works, these data need to be created with the design from God. If the Will of God exists, then it must be in these parameters - seed+tags. These params finally create picture inside artifficial intelligence.",
    "three": "True-random numbers, quantum random numbers are used to generate seed+tags input data. Using mathematical transformations and dictionaries, a Seed and a set of tags are formed from them. (all the Artchitect's code is always available on ",
    "four": "Quantum random numbers are obtained by reading the quantum background. To do this, used webcam and normalized the whole frame to one 64-bit number. All 64 bits are obtained from the webcam matrix, on which light quanta fall. (if someone is interested in more details, write me an email and I will tell you in detail)",
    "five": "Quants of light are created in the surrounding Universe constantly. The meaning of the frame itself does not matter to the system, but I don't know this exactly. I have, for example, Jesus there :) I think to me that the ritual component is also important here.",
    "example": "Example of my frame from which quantum noise is being read:",
    "remaining": "It remains only to get meaningful beam of light.",
    "chain": "The whole chain of events is as follows:",
    "remaining_line1": "meaningful light-beam being converted into meaningful digital data inside a webcam",
    "remaining_line2": "meaningful digital data being converted into meaningful numbers",
    "remaining_line3": "meaningful numbers being converted into meaningful seed and tags",
    "remaining_line4": "meaningful seed+tags being converted into a meaningful picture inside artificial intelligence (which is just a digital artist).",
    "meaningful_quants": "Meaningful quants of light can create only God, if He will. God can tune lights to reach final goal, if He want to.",
    "fill_sense": "If God wants to fill Artchitect with some sense, He have a open-API to do it.",
    "sinse": "Since this text would not have been written without the Will of God, He also knows this way, and has always known."
  },
  "ru": {
    "how": "Как это технически устроено",
    "zero": " На выходе мы хотим получить миллионы художественных картин, созданных по Божьему Замыслу. Мы попробуем увидеть этот замысел, если Он есть. Или же не увидим, если Его нет.",
    "one": "Сама картина технически рисуется современным искусственным интеллектом, в котором можно создать бесконечное(это важно) количество творений. В Artchitect используется Stable Diffusion. Бог может выбрать любой вариант картины, какой захочет. Он знает все возможные варианты заранее.",
    "two": "Для формирования картины в Stable Diffusion нужно указать набор текстовых тегов и специальное число - Seed. Эти данные нужно сформировать, учитывая Волю Бога. Если Воля Бога существует, то она должна быть в этих самых параметрах seed+tags, так как именно через них формируется уникальная картина в нейронной сети.",
    "three": "Для формирования входных данных seed+tags используются true-случайные числа, квантовые случайные числа. Из них с помощью математических преобразований и словарей формируется Seed и набор тегов. (весь код архитектора всегда доступен на",
    "four": "Квантовые случайные числа получаются через считывание квантового фона. Для этого используетс веб-камера  нормализация целого кадра до одного 64-битного числа. Все 64 бита получаются с матрицы веб-камеры, на котору падают кванты света. (если кому-то интересно более подробно, напишите мне на почту и я расскажу в деталях)",
    "five": "Кванты света формируются в окружающей Вселенной постоянно, и сама камера формирует эти лучи (там есть диоды на ней). Имеет ли смысл само содержимое кадра - я не знаю.",
    "example": "Фото моего кадра, с которого считывается квантовый шум (у меня там икона Иисуса Христа):",
    "remaining": "Осталось лишь получить осмысленные кванты света.",
    "chain": "Вся цепочка событий такова:",
    "remaining_line1": "осмысленные кванты света превращаются в осмысленные цифровые данные внутри веб-камеры",
    "remaining_line2": "осмысленные цифровые данные превращаются в осмысленные числа",
    "remaining_line3": "осмысленные числа превращаются в осмысленные seed и теги",
    "remaining_line4": "осмысленные seed+tags превратятся в осмысленную картину внутри искусственного интеллекта (который тут просто цифровой художник, исполнитель)",
    "meaningful_quants": "Осмысленные кванты света может сформировать только Бог, если Он так пожелает. В нужный момент времени Бог может так подыграть освещение, что камера поймает осмысленный луч света.",
    "fill_sense": "Если Бог захочет наполнить смыслом Архитектор, то у Него есть способ это сделать.",
    "sinse": "Так как и этот текст не был бы написан без Воли Бога, то Он тоже знает этот способ, и всегда знал. Есть вероятность, что работа Архитектора вообще будет зависеть от текста, написанного на этой странице."
  }
}
</i18n><template>
  <div class="mt-5">
    <h1 class="is-size-4">{{$t('how')}}</h1>

    <p><b>0)</b> {{$t('zero')}}</p>

    <p><b>1)</b> {{$t('one')}}</p>

    <p><b>2)</b> {{$t('two')}}</p>

    <p><b>3)</b> {{$t('three')}} <a href="https://github.com/artchitector/artchitect"
                     target="_blank">github.com/artchitector/artchitect</a>)</p>

    <p><b>4)</b> {{$t('four')}}</p>

    <p><b>5)</b> {{$t('five')}}<br/>{{$t('example')}}
    </p>

    <p>
      {{$t('remaining')}}<br/>
      {{$t('chain')}}<br/>
      - {{$t('remaining_line1')}}<br/>
      - {{$t('remaining_line2')}}<br/>
      - {{$t('remaining_line3')}}<br/>
      - {{$t('remaining_line4')}}<br/>
    </p>

    <p>{{$t('meaningful_quants')}}</p>

    <p>{{$t('fill_sense')}}</p>

    <p>{{$t('sinse')}}</p>
  </div>
</template>
