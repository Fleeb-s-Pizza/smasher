const { createApp } = Vue;

const app = createApp({
    data() {
        return {
            tabs: [
                'Images',
                'Documents',
                'Info'
            ],
            currentTab: 'Images',
            images: {
                webp: false,
                url: '',
                width: 0,
                height: 0,
                quality: 100,
                rotate: 0,
            },
            info: null
        }
    },
    methods: {
        getSmashedUrl: function(images) {
            return "https://smasher.fleebs.gg/image?url="+ encodeURI(images.url) + (images.width > 0 ? '&width=' + images.width : '') + (images.height > 0 ? '&height=' + images.height : '') + (images.webp ? '&webp=true' : '') + (images.quality < 100 ? '&q=' + images.quality : '') + (images.rotate > 0 ? '&rotate=' + images.rotate : '');
        }
    },
    mounted() {
        axios.get('/info')
            .then((response) => {
                console.log(response.data);
                this.info = response.data;
                console.log(this.info);
            })
    }
});

app.mount('#app');