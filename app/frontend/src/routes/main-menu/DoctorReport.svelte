<script lang="ts">
    import { doctorReportStore, settingsStore } from "../../models/stores";
    
</script>

{#if !$doctorReportStore}
    <div class="container d-flex align-items-center justify-content-center">
        <div class="text-center">
            <div class="spinner-border text-primary mb-3" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
            <h4 class="text-muted">Checking Required Dependencies...</h4>
        </div>
    </div>
{:else}
    <div class="container mt-5">
        <div class="row justify-content-center">
            <div class="col-md-12">
                <div class="card">
                    <div class="card-header bg-dark text-white">
                        <h3 class="card-title mb-0">Required Dependencies</h3>
                    </div>
                    <div class="card-body">
                        <ul class="list-group list-group-flush">
                            <li class="list-group-item d-flex justify-content-between align-items-center">
                                FFMPEG
                                <span class="badge {$doctorReportStore.SystemHasFFMPEG ? 'bg-success' : 'bg-danger'} rounded-pill">
                                    {$doctorReportStore.SystemHasFFMPEG ? "INSTALLED" : "MISSING"}
                                </span>
                            </li>
                            <li class="list-group-item d-flex justify-content-between align-items-center">
                                ARW Support
                                <span class="badge {$doctorReportStore.SystemHasARWConversion ? 'bg-success' : 'bg-danger'} rounded-pill">
                                    {$doctorReportStore.SystemHasARWConversion ? "INSTALLED" : "MISSING"}
                                </span>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>
{/if}
