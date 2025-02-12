   <div class="container height clean">
    <div class="row height">
        
      <div class="col-md-3 border-right height pt-4 accent-color d-flex flex-column align-items-center">
        <a href="#/main-menu" class="btn btn-secondary btn-block main-color my-buttons">
            Quit
        </a>
        <a href="#/main-menu" class="btn btn-secondary btn-block secondary-color my-buttons mt-2">
            Save and Quit
        </a>
      </div>
      <div class="col-md-9 p-3">
        <h2>Astropath Settings</h2>
        <div class="form-check form-switch">
            <input class="form-check-input" type="checkbox" role="switch" id="flexSwitchCheckDefault">
            <label class="form-check-label" for="flexSwitchCheckDefault">Enable Time Lapse Generation</label>
            <p class="description">Brief description of feature goes here</p>
        </div>
        <div class="form-check form-switch">
            <input class="form-check-input" type="checkbox" role="switch" id="flexSwitchCheckDefault">
            <label class="form-check-label" for="flexSwitchCheckDefault">Support ARW Conversions</label>
            <p class="description">Brief description of feature goes here</p>
        </div>
      </div>
    </div>
  </div>
  
  <style>
    .border-right {
      border-right: 1px solid rgb(177, 177, 177);
    }
    .height{
        height: 100%;
    }
    
    .main-color{
        background-color: #1E5A7E;
    }
    .secondary-color{
        background-color: #7E1E5A ;
    }
    .accent-color{
        background-color: #F0EDE5;
    }
    .my-buttons{
        width: 80%;
    }
    .clean{
        margin-left: 0%;
        padding-left: 0;
    }
    .description{
        font-size: smaller;
        color: #8f8f8f;
    }
  </style>
