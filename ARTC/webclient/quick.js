function hex(buffer) {
  const hexCodes = [];
  const view = new DataView(buffer);
  for (let i = 0; i < view.byteLength; i += 4) {
    const value = view.getUint32(i);
    const stringValue = value.toString(16);
    const padding = '00000000';
    const paddedValue = (padding + stringValue).slice(-padding.length);
    hexCodes.push(paddedValue);
  }
  return hexCodes.join('');
}
async function generateToken(appId, appKey, channelId, userId, timestamp) {
  const encoder = new TextEncoder();
  const data = encoder.encode(`${appId}${appKey}${channelId}${userId}${timestamp}`);

  const hash = await crypto.subtle.digest('SHA-256', data);
  return hex(hash);
}

function showToast(baseId, message) {
  $(`#${baseId}Body`).text(message);
  const toast = new bootstrap.Toast($(`#${baseId}`));

  toast.show();
}

// 填入您的应用ID 和 AppKey
const appId = '4855ee4e-dab9-4694-b392-e8f38bd48dbd';
const appKey = '415ebb6459519642da7bd1292e0c11fb';
AliRtcEngine.setLogLevel(0);
let aliRtcEngine;
const remoteVideoElMap = {};
const remoteVideoContainer = document.querySelector('#remoteVideoContainer');

function removeRemoteVideo(userId, type = 'camera') {
  const vid = `${type}_${userId}`;
  const el = remoteVideoElMap[vid];
  if (el) {
    aliRtcEngine.setRemoteViewConfig(null, userId, type === 'camera' ? 1: 2);
    el.pause();
    remoteVideoContainer.removeChild(el);
    delete remoteVideoElMap[vid];
  }
}

function listenEvents() {
  if (!aliRtcEngine) {
    return;
  }
  // 监听远端用户上线
  aliRtcEngine.on('remoteUserOnLineNotify', (userId, elapsed) => {
    console.log(`用户 ${userId} 加入频道，耗时 ${elapsed} 秒`);
    // 这里处理您的业务逻辑，如展示这个用户的模块
    showToast('onlineToast', `用户 ${userId} 上线`);
  });

  // 监听远端用户下线
  aliRtcEngine.on('remoteUserOffLineNotify', (userId, reason) => {
    // reason 为原因码，具体含义请查看 API 文档
    console.log(`用户 ${userId} 离开频道，原因码: ${reason}`);
    // 这里处理您的业务逻辑，如销毁这个用户的模块
    showToast('offlineToast', `用户 ${userId} 下线`);
    removeRemoteVideo(userId, 'camera');
    removeRemoteVideo(userId, 'screen');
  });

  aliRtcEngine.on('bye', code => {
    // code 为原因码，具体含义请查看 API 文档
    console.log(`bye, code=${code}`);
    // 这里做您的处理业务，如退出通话页面等
    showToast('loginToast', `您已离开频道，原因码: ${code}`);
  });

  aliRtcEngine.on('videoSubscribeStateChanged', (userId, oldState, newState, interval, channelId) => {
    // oldState、newState 类型均为AliRtcSubscribeState，值包含 0（初始化）、1（未订阅）、2（订阅中）、3（已订阅）
    // interval 为两个状态之间的变化时间间隔，单位毫秒
    console.log(`频道 ${channelId} 远端用户 ${userId} 订阅状态由 ${oldState} 变为 ${newState}`);
    const vid = `camera_${userId}`;
    // 处理示例
    if (newState === 3) {
      const video = document.createElement('video');
      video.autoplay = true;
      video.className = 'video';
      remoteVideoElMap[vid] = video;
      remoteVideoContainer.appendChild(video);
      // 第一个参数传入 HTMLVideoElement
      // 第二个参数传入远端用户 ID
      // 第三个参数支持传入 1 （预览相机流）、2（预览屏幕共享流）
      aliRtcEngine.setRemoteViewConfig(video, userId, 1);
    } else if (newState === 1) {
      removeRemoteVideo(userId, 'camera');
    }
  });

  aliRtcEngine.on('screenShareSubscribeStateChanged', (userId, oldState, newState, interval, channelId) => {
    // oldState、newState 类型均为AliRtcSubscribeState，值包含 0（初始化）、1（未订阅）、2（订阅中）、3（已订阅）
    // interval 为两个状态之间的变化时间间隔，单位毫秒
    console.log(`频道 ${channelId} 远端用户 ${userId} 屏幕流的订阅状态由 ${oldState} 变为 ${newState}`);
    const vid = `screen_${userId}`;
    // 处理示例    
    if (newState === 3) {
      const video = document.createElement('video');
      video.autoplay = true;
      video.className = 'video';
      remoteVideoElMap[vid] = video;
      remoteVideoContainer.appendChild(video);
      // 第一个参数传入 HTMLVideoElement
      // 第二个参数传入远端用户 ID
      // 第三个参数支持传入 1 （预览相机流）、2（预览屏幕共享流）
      aliRtcEngine.setRemoteViewConfig(video, userId, 2);
    } else if (newState === 1) {
      removeRemoteVideo(userId, 'screen');
    }
  });
}

$('#loginForm').submit(async e => {
  // 防止表单默认提交动作
  e.preventDefault();
  const channelId = $('#channelId').val();
  const userId = $('#userId').val();
  const timestamp = Math.floor(Date.now() / 1000) + 3600 * 3;

  if (!channelId || !userId) {
    showToast('loginToast', '数据不完整');
    return;
  }

  aliRtcEngine = AliRtcEngine.getInstance();
  listenEvents();

  try {
    const token = await generateToken(appId, appKey, channelId, userId, timestamp);
    // 设置频道模式，支持传入字符串 communication（通话模式）、interactive_live（互动模式）
    aliRtcEngine.setChannelProfile('communication');
    // 设置角色，互动模式时调用才生效
    // 支持传入字符串 interactive（互动角色，允许推拉流）、live（观众角色，仅允许拉流）
    // aliRtcEngine.setClientRole('interactive');
    // 加入频道，参数 token、nonce 等一般有服务端返回
    await aliRtcEngine.joinChannel(
      {
        channelId,
        userId,
        appId,
        token,
        timestamp,
      },
      userId
    );
    showToast('loginToast', '加入频道成功');
    $('#joinBtn').prop('disabled', true);
    $('#leaveBtn').prop('disabled', false);

    // 预览
    aliRtcEngine.setLocalViewConfig('localPreviewer', 1);
  } catch (error) {
    console.log('加入频道失败', error);
    showToast('loginToast', '加入频道失败');
  }
});

$('#leaveBtn').click(async () => {
  Object.keys(remoteVideoElMap).forEach(vid => {
    const arr = vid.split('_');
    removeRemoteVideo(arr[1], arr[0]);
  });
  // 停止本地预览
  await aliRtcEngine.stopPreview();
  // 离开频道
  await aliRtcEngine.leaveChannel();
  // 销毁实例
  aliRtcEngine.destroy();
  aliRtcEngine = undefined;
  $('#joinBtn').prop('disabled', false);
  $('#leaveBtn').prop('disabled', true);
  showToast('loginToast', '已离开频道');
});