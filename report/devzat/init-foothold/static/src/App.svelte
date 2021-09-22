<script>
  import { onMount } from 'svelte';

  let animals = [];
  let postAnimal = { name: '', species: 'cat' };

  onMount(async () => {
    // fetch animals
    const res = await fetch(`/api/pet`);
    animals = await res.json();
  });

  async function doPost() {
    // add the new animal
    await fetch(`/api/pet`, {
      method: 'POST',
      body: JSON.stringify(postAnimal),
    })
      .then(async (res) => {
        if (res.status == 200) {
          // reload animals
          const update = await fetch(`/api/pet`);
          animals = await update.json();

          // Clear form
          postAnimal = { name: '' };
        }
      })
      .catch((err) => alert(err));
  }
</script>

<main class="container">
  <div class="row mb-5">
    <h1 align="center">Pet Inventory</h1>

    <p align="center">
      Welcome to my pet inventory. This is where I keep a list of my pets.
    </p>
    <p align="center">I mean, come one, who doesn't like animals, right?</p>
  </div>

  <div class="row mb-5">
    <h2 align="center">My Pets</h2>
    <table class="table table-striped table-hover">
      <thead>
        <tr>
          <th scope="col">Name</th>
          <th scope="col">Species</th>
          <th scope="col">Characteristics</th>
          <th scope="col" />
        </tr>
      </thead>
      <tbody>
        {#each animals as { name, species, characteristics }, i}
          <tr id={i}>
            <td>{name}</td>
            <td>{species.charAt(0).toUpperCase() + species.slice(1)}</td>
            <td>{characteristics}</td>
            <td
              ><button
                class="btn btn-outline-secondary"
                on:click|preventDefault={() => alert('Not implemented, yet')}
                ><i class="fas fa-trash" /></button
              ></td
            >
          </tr>
        {/each}
      </tbody>
    </table>
  </div>

  <div class="row mb-5">
    <h2 align="center">Add a Pet</h2>
    <form>
      <div class="mb-3">
        <label for="name" class="form-label">Name the pet</label>
        <input
          type="text"
          class="form-control"
          name="name"
          id="name"
          bind:value={postAnimal.name}
        />
      </div>

      <div class="mb-3">
        <label for="species" class="form-label">Which species is it?</label>
        <select
          class="form-select"
          name="species"
          id="species"
          bind:value={postAnimal.species}
        >
          <option value="cat">Cat</option>
          <option value="dog">Dog</option>
          <option value="gopher">Gopher</option>
          <option value="giraffe">Giraffe</option>
          <option value="redkite">Red Kite</option>
          <option value="bluewhale">Blue Whale</option>
        </select>
      </div>

      <button
        class="btn btn-outline-secondary"
        type="submit"
        on:click|preventDefault={doPost}>Add Pet</button
      >
    </form>
  </div>
</main>

<style>
</style>
