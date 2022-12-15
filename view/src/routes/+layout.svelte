<script>

</script>

<svelte:head>
    <title>artchitect.space</title>
</svelte:head>

<main>
    <header>
        <h1><a href="https://artchitect.space">artchitect.space</a></h1>
    </header>

    <slot></slot>

    <footer>
        github: <a href="https://github.com/artchitector/artchitect" target="_blank">https://github.com/artchitector/artchitect</a>
    </footer>
</main>

<style>
    body {
        padding: 0;
        margin: 0 0 100px 0;
    }

    header {
        margin: 0;
        background-color: #2d2d2d;
        text-align: center;
        padding: 10px 0 0 0;
    }

    header h1 {
        margin: 0;
        padding: 0;
        color: #d4d1c3;
    }

    header h1 a {
        color: #d4d1c3;
        text-decoration: none;
    }

    footer {
        background-color: #2d2d2d;
        color: #f1f1f1;
        text-align: center;
        padding: 0 0;
    }
    footer a {
        color: #d4d1c3;
        text-decoration: none;
    }
</style>