<template>
    <div class="container">
        <div v-if="!comment.clicked" class="reaction">
            <button @click="toggle('up')" :class="firstClass" class="reaction-div tooltip like-button">
                <font-awesome-icon :icon="['fas', 'thumbs-up']" />
                <span class="tooltiptext tooltiptext-position-left">Breaking news.</span>
            </button>
            <button @click="toggle('down')" :class="secondClass" class="reaction-div tooltip dislike-button">
                <font-awesome-icon :icon="['fas', 'thumbs-down']" />
                <span class="tooltiptext tooltiptext-position-right">Not breaking news.</span>
            </button>
            <div @click="toggleFeedbackForm('show-feedback')" v-show="clicked.up || clicked.down" class="comment-container">
                <span class="comment-action" >Comment</span>
            </div>
        </div>
        <div v-else-if="comment.clicked && !feedBackSent" :class="comment.class" class="feedback">
            <input v-model="comment.text" placeholder="Your feedback ..." class="feedback-input" ref="inputRef" type="text" />
            <div class="send-or-cancel">
                <span @click="toggleFeedbackForm('')">Cancel</span>
                <button class="feedback-submit" @click="submitForm()">Send</button>
            </div>
        </div>
        <div v-else class="feedback-gratitude">
            <span>Thank you for your feedback!</span>
        </div>
    </div>
</template>
<script>
    import { library } from '@fortawesome/fontawesome-svg-core';
    import { faThumbsUp, faThumbsDown} from '@fortawesome/free-solid-svg-icons';
    import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
    import './Reaction.css';
    import { post } from "../../utils/api";

    library.add(faThumbsUp, faThumbsDown);
    export default {
        name: 'ReactionLike',
        props: {
            article: Object
        },
        components: {
            FontAwesomeIcon
        },
        data() {
            return {
                clicked: {
                    up: false,
                    down: false
                },
                firstClass: '',
                secondClass: '',
                comment: {
                    text: '',
                    clicked: false,
                    class: 'show-feedback'
                },
                feedBackSent: false
            }
        },
        methods: {
            async toggle(which) {
                const article = this.article;
                if(localStorage.getItem(article["id"]) !== null){
                    return;
                }
                if(which === 'up'){
                    this.clicked.up = true;
                } else if(which === 'down'){
                    this.clicked.down = true;
                }
                localStorage.setItem(article["id"], true);
                await post('rate', null, {
                    article: this.article,
                    rating: this.clicked.up ? "like" : "dislike",
                });
            },
            async toggleFeedbackForm(className) {
                this.comment.clicked = !this.comment.clicked;
                this.comment.class = className;
            },
            async submitForm() {
                this.comment.clicked = true;
                this.comment.class = 'feedback-gratitude';
                this.feedBackSent = true;
                await post('comment', null, {
                    id: this.article["id"],
                    project: this.article["project"],
                    comment: this.comment.text,
                });
            },
        },
        watch: {
            clicked: {
                handler(newValue) {
                    if(newValue.up){
                        this.firstClass = 'like-clicked';
                        this.secondClass = 'disabled';
                    } else if(newValue.down){
                        this.firstClass = 'disabled';
                        this.secondClass = 'dislike-clicked';
                    }
                },
                deep: true
            },
            comment: {
                handler(newValue) {
                    if(newValue.clicked){
                        this.comment.clicked = true;
                    }
                },
                deep: true
            }
        }
    }
</script>