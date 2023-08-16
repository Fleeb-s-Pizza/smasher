const { createApp } = Vue;

const app = createApp({
    data() {
        return {
            tabs: [
                'Images',
                'Documents',
            ],
            currentTab: 'Images',
            images: {
                webp: false,
                url: '',
                width: 0,
                height: 0,
            }
        }
    },
    methods: {
        getSmashedUrl: function(images) {
            console.log(encodeURI(images.url));
            return "https://smasher.fleebs.gg/image?url="+ encodeURI(images.url) + (images.width > 0 ? '&width=' + images.width : '') + (images.height > 0 ? '&height=' + images.height : '') + (images.webp ? '&webp=true' : '');
        }
    }
});

app.mount('#app');