function formatTime(seconds) {
    seconds = Math.max(0, seconds);

    const hours = Math.floor(seconds / 3600);

    const minutes = Math.floor(
        (seconds % 3600) / 60
    );

    const secs = seconds % 60;

    return (
        String(hours).padStart(2, '0') +
        ':' +
        String(minutes).padStart(2, '0') +
        ':' +
        String(secs).padStart(2, '0')
    );
}


function startCooldown(seconds) {

    const timer =
        document.getElementById('timer');

    if (!timer) {
        return;
    }


    function update() {

        timer.textContent =
            formatTime(seconds);


        if (seconds <= 0) {

            clearInterval(interval);

            timer.textContent =
                '00:00:00';

            return;
        }


        seconds--;
    }


    update();


    const interval =
        setInterval(update, 1000);
}


async function copyKey() {

    const keyElement =
        document.getElementById('promo-key');

    const button =
        document.getElementById('copy-button');


    if (!keyElement || !button) {
        return;
    }


    const key =
        keyElement.textContent.trim();


    try {

        await navigator.clipboard.writeText(key);

        const oldText =
            button.textContent;

        button.textContent =
            '✓';

        setTimeout(() => {
            button.textContent =
                oldText;
        }, 1500);

    } catch (error) {

        /*
         * Fallback для старых браузеров.
         */

        const textarea =
            document.createElement('textarea');

        textarea.value = key;

        document.body.appendChild(
            textarea
        );

        textarea.select();

        document.execCommand('copy');

        textarea.remove();


        button.textContent = '✓';

        setTimeout(() => {
            button.textContent = '📋';
        }, 1500);
    }
}