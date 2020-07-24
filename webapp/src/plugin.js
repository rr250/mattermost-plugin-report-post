/* eslint-disable lines-around-comment */
/* eslint-disable no-console */
import {Component} from 'react';
import {v4 as uuidv4} from 'uuid';
import axios from 'axios';

import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {getCurrentUser, getUser} from 'mattermost-redux/selectors/entities/users';

import {id as pluginId} from './manifest';

// eslint-disable-next-line react/require-optimization
class ReportPlugin extends Component {
    initialize(registry, store) {
        registry.registerPostDropdownMenuAction('Report', (postId) => {
            const state = store.getState();
            const post = getPost(state, postId);
            const currentUser = getCurrentUser(state);
            const user = getUser(state, post.user_id);
            axios.post('/plugins/' + pluginId + '/postreport', {
                id: uuidv4(),
                reported_by: currentUser.first_name + ' ' + currentUser.last_name,
                reported_by_id: currentUser.id,
                created_at: new Date().toISOString(),
                reported_name: user.first_name + ' ' + user.last_name,
                reported_id: user.id,
                channel_id: post.channel_id,
                reported_username: user.username,
                reported_email: user.email,
                reported_text: post.message,
                reported_text_id: post.id,
            }).then().catch((err) => {
                console.log(err);
            });
        });
    }

    uninitialize() {
        console.log(pluginId + '::uninitialize()');
    }
}

export default ReportPlugin;
