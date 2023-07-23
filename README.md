# Artchitect

![artchitect_logo](https://github.com/artchitector/artchitect/blob/master/eye/static/jesus_anim_92.gif)

#### https://artchitect.space

> Artchitect - it is an amazing autonomous creative machine capable of creating magnificent artworks inspired by the
> universe around us. In its continuous creativity , the machine receives inspiration from the natural entropy of the
> universe, represented as background light, and creates unique artworks without human participation.

Techically artchitect-project is the control-system wrapped around an art-system - Stable Diffusion AI v1.5. Stable
Diffusion AI is the ability to draw pictures for Artchitect.
The data to run Architect is light background/noise. The light background is read using a webcam, the frame is converted
to an int64 number. **Int64** is the source of all solutions (randomly select something from the list or create a unique
initial value).

### architecture of artchitect

Two parts:

- soul - home computer with strong GPU to run Stable
  Diffusion (artchitect using RTX 3060 12Gb and Stable Diffusion v1.5, Invoke.AI "fork").
- gate - multiple dedicated VDSs (frontend+backend server, file storage, database)

There is no reliability question, SLA is not important. If the Artchitect breaks down and turns off, it can stand for
several days
before repair. VDS servers that provide viewing of paintings work reliably, but the "soul" - main home computer can
be turned off for downtime from days to weeks
for downtime, up to several weeks, if necessary.

Home computer need access to VDS, but VDS doesn't need access to home computed.

More GPU-RAM = larger resolution = more quality. RTX3060-12Gb gives resolution 2560x3840 (printable on 40x60 canvas).

golang backend services + python backend services, splitted between home computer and remote VDS (visible from
Internet).

### manual

🤝 there are no instructions for Artchitect, since no one needed it before.
If you need more information or instructions to install your copy of Artchitect - please ask me questions in issues. I
will help.

### engineering style

- artchitect is available for everyone. If you want your running-copy of Artchitect, you need devops skills (to
  setup and run your own servers, databases, to understand code)

- artchitect is one-person-project, so there is no code style here

- many things in source-code are not called obviously: main
  service called "soul", and it consists of "gifter", "merciful", "speller"... (sometimes the author does not remember
  which of them does what)

- no docs. if you need help, make an issue. author will personally explain what to do or will make docs

- no braches, everything is in master in singlerepo. Almost every commit is delivered to production.

- no testing, neither manual nor automated

- manual deployment with simple ssh-commands, crontasks

### How Artchitect looks like:

![artchitect_installation](https://github.com/artchitector/artchitect/blob/master/eye/static/artchitect_in_real_world.jpg)