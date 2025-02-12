<script lang="ts">

    import DoctorReport from "./DoctorReport.svelte";
    import { doctorReportStore } from "../../models/stores";
    
</script>

<section>
    <div class="row justify-content-center align-items-center" style="height: 100vh;">
        <div class="col-md-8 text-center">
            <img src="src/assets/images/ASTROPATH-removebg-preview.png" alt="Astropath logo" class="mb-4" style="max-width: 500px;" />
            <div class="d-grid gap-3 mt-4">
                <a href="#/upload" class="btn btn-primary btn-lg w-100 main-color">
                    Get Started
                </a>
                <div class="d-flex gap-2">
                    <a href="#/settings" class="btn btn-secondary btn-lg flex-fill secondary-color">
                        Settings
                    </a>
                    <a href="#/info" class="btn btn-danger btn-lg flex-fill secondary-color">
                        Info
                    </a>
                </div>
            </div>
            <div class="mt-5">
                <DoctorReport />
            </div>
        </div>
    </div>
</section>

<style>
    .main-color{
        background-color: #1E5A7E;
    }
    .secondary-color{
        background-color: #7E1E5A ;
    }
</style>
  
