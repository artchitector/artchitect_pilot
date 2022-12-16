<script>
    import {onMount} from 'svelte';

    let state = {
        loading: false,
        state: null,
        error: null,
    };

    onMount(async () => {
        // TODO Need do 5sec auto-reloading
        await loadState();
    })

    async function loadState() {
        try {
            state.error = null;
            state.loading = true;
            const response = await fetch(import.meta.env.VITE_API_HOST + "/state");
            const json = await response.json();
            console.log("State received: ", json)
            state.state = json;

        } catch (err) {
            state.error = err;
        } finally {
            state.loading = false;
        }
    }
</script>

<div>
    <button on:click={() => loadState()}>reload</button>
    {#if state.loading}
        <p>...loading</p>
    {:else if state.error}
        <p class="danger">{state.error}</p>
    {:else if state.state}
        <pre><b>Current state</b>: {JSON.stringify(state.state)}</pre>
        <hr/>
        <div class="lastimage_block">
            <pre><b>Current image (id={state.state.State.LastPainting.ID}): "{state.state.State.LastPainting.Caption}"</b></pre>
            <img class="lastimage" src="{import.meta.env.VITE_API_HOST}/painting/{state.state.State.LastPainting.ID}"/>
        </div>
    {/if}
</div>

<style>
    div.lastimage_block {
        text-align: center;
    }
    img.lastimage {
        width: 100%;
        max-width: 300px;
    }
    p.danger {
        color: red;
    }
</style>