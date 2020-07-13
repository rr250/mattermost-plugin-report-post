/* eslint-disable no-console */
import {Component} from 'react';
import {v4 as uuidv4} from 'uuid';

import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {getCurrentUser, getUser} from 'mattermost-redux/selectors/entities/users';
import {Client4} from 'mattermost-redux/client';

import {id as pluginId} from './manifest';

// eslint-disable-next-line react/require-optimization
class ReportPlugin extends Component {
    initialize(registry, store) {
        registry.registerPostDropdownMenuAction('Report', (postId) => {
            const state = store.getState();
            const post = getPost(state, postId);
            const currentUser = getCurrentUser(state);
            const user = getUser(state, post.user_id);
            Client4.getConfig().then((res) => {
                const botId = res.PluginSettings.Plugins[pluginId].botid;
                const channelId = res.PluginSettings.Plugins[pluginId].channelid;
                const newPost = {
                    pending_post_id: uuidv4(),
                    user_id: botId,
                    channel_id: channelId,
                    message: '',
                    props: {
                        attachments: [{
                            text: 'Report Alert:-\n\tReported: ' +
                            user.first_name + ' ' + user.last_name +
                            '\n\tReported ID: ' + user.id +
                            '\n\tReported Channel ID: ' + post.channel_id +
                            '\n\tReported Username: ' + user.username +
                            '\n\tReported Email: ' + user.email +
                            '\n\tReported By: ' + currentUser.username +
                            '\n\tReported By ID: ' + currentUser.id +
                            '\n\tReported Text ID: ' + post.id +
                            '\n\nReported Text:-\n' + post.message + '\n',
                        }],
                    },
                };
                Client4.createPost(newPost);
            });
        });
    }

    uninitialize() {
        console.log(pluginId + '::uninitialize()');
    }
}

export default ReportPlugin;
