function fetchImages() {
    return Promise.resolve([
      {
        Title: "Sunset",
        AltText: "Clouds at sunset",
        URL: "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
      },
      {
        Title: "Mountain",
        AltText: "A mountain at sunset",
        URL: "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
      },
    ]);
  }

  function timeout(t, v) {
    return new Promise(res => {
      setTimeout(() => res(v), t);
    })
  }

  const gallery$ = document.querySelector(".gallery");

  timeout(2000, fetchImages()).then(images => {
    gallery$.textContent = images.length ? "" : "No images available.";

    images.forEach(img => {
      const imgElem$ = document.createElement("img");
      imgElem$.src = img.URL;
      imgElem$.alt = img.AltText;
      const titleElem$ = document.createElement("h3");
      titleElem$.textContent = img.Title;
      const wrapperElem$ = document.createElement("div");
      wrapperElem$.classList.add("gallery-image");
      wrapperElem$.appendChild(titleElem$);
      wrapperElem$.appendChild(imgElem$);
      gallery$.appendChild(wrapperElem$);
    });
  });