var timeoutInterval = 1000;
var refreshIntervalID = 0;
const statusEndpoint = "/status"

document.addEventListener("DOMContentLoaded", function () {
    const gaugeContainer = document.getElementById("gauge");
    // const statusValue = document.getElementById("status").dataset.teststatus;

    // Load the Lottie animation
    const gaugeAnimation = lottie.loadAnimation({
        container: gaugeContainer,
        renderer: 'svg',
        loop: false,
        autoplay: false,
        path: 'static/animations/gauge.json' // path to your JSON
    });



    function refresh() {
        gaugeAnimation.goToAndPlay(50, true);
        fetch(window.location.origin + statusEndpoint)
            .then((response) => {
                if (response.ok) {
                    gaugeAnimation.goToAndStop(50, true);
                    showResults();
                    window.clearInterval(refreshIntervalID);
                    window.location.replace(window.location.href);
                }
            })
            .catch(error => {
                console.error(error);
            });
    }

    window.addEventListener("load", (event) => {
        if (document.getElementById("testButton").disabled) {
            refreshIntervalID = setInterval(refresh, timeoutInterval);
        }
    });
    showResults();
    function showResults() {
        document.querySelectorAll(".result-card").forEach((card, index) => {
            setTimeout(() => card.classList.add("show"), index * 200);
        });
    }
});
