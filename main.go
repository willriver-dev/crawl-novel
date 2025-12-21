package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// CONFIGURATION
// ============================================================================

const (
	GraphQLEndpoint = "https://bff-page.kakao.com/graphql"

	// Thay Cookie mới của bạn vào đây
	CookieString = "_kpdid=fe3277a6e854465f872802cdf1e5891c; _kpiid=278f88bc29a85822679573a2aa34ec0f; _kau=a710fc48b7494608758b095aded223cb596771d398c91e2d92d5c433e345fa34ecb776fee5e479588bd124266dc8a5f0077e966b6c14d681d8fddf7cfc65a8a027c45df8d38636ea39a9dbec4952d06ac22664ace960a1773548c96417089ab01d0629e584cd690a443cb43c819b43a3c964043f6a8ba7c80c49f76b36373830323833393532303039393832363331343639393639373536363832388e89d40926e528c875cd450a1de2119c; _kawlt=eEj-yThNzqPPZY1sp8SXm-5KqIKARSyh3ytXMQNRChq_U2jJikGjhXG4O5uMr7Hnvalj8-WFL_mffT8Kc_lPKru71w0VN2SvDCEL7uVgK0lX8YgKlZtb9iwUM272aBP5; _kawltea=1766154377; _karmt=tzyD4XIqRdHS3Laklr6xnAmr4rf1R3ka_nCe3pMp7a6MVK0rc1L7d2Uwe64LqcAB; _karmtea=1768681577; _kahai=ba4be1e91c915262b60a6402f5e1e321446b1a1739f6a79ad21e84f1393c8927; _kaslt=nu7qfiHmRXQp/QXcXBcfJfkm9myUu8MEFefr9zaMrvAiAWWuPQy+MKs86M6YJmhCsOGZ5qTkO+PaOmsHd187gQ==; _kpawbat_e=r0L30lJDjZnaKz6F5zhaGME83e7j3%2BrqJPR8NfMg2Sv454BQn4y8lQe1r9FWX15gz%2Boc%2F1MdvA7oN4b1dVy6ahybSF2a5etY2NQ56dhbX8M%3D; _kpwtkn=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..l5TAO0gCUMXsRmh2.i_W_W8fDz3wP0WdGb7oTp1e7F84AviJ3hKzd_TJT0rUm8qEYwPHzRlW7i3A_kSIpLTRILhHJM_stsh8-KtIvSbQhCKlNzB99EoHQYmnXuFwYYzV3OkySfe4J7NT82_nW0XCdJbFjpYDyPW7GTPci1DfRrGnwuDgierJ_2vE4_kDwpZRHhI1xF89VeutU76gu82YpQqteoEUjjucU6YOYA5yBbN1dh1YmQZ7MpSTTsXGWvBaFnbb50ZHt-sZlBVcmYDgiXodKGws3itfNsedhGhOGE7mUDhX4X2l29hwcRP3ZhXacfLNrPDmJEH52x1WSFrMdcHFvi4okDgpp1NGF_EnBzOwYlTkVos47Nt03npexHBr8R4GxT6fmZt74FOlodu61SZSsFm_OdK9teqizCYD0LZcmvrrdB3Qea5fiQXj2x1Us2dKCvF2tN66BWPIIdPhg5zssll2VqY3_ggkA2qu8y4jB96s84GLHo-GycAX6YD1jVTVCWWteXfNejd0ox2dZPd6K6x1KH6b8Z1-670-Mj8OJZ15UdWRCf5T5rtKKTLLTdHdqUnLowWgGzlVKrZAXLzJ-v4-po-ZrW-LIy7SMsx_CPEyb_jMW2a2Ync2rxzw-Zc5F0UY1BN_2DxKuLTa_70KUC1DRRB-1tEZR0lzRrvw6B2NzO6Ct3PQVS2mr2JrGbZofNqyjx7FAH-KwShTruCntLS1HHVYd4U5PqwedF940GKy-MwZ6e-RHfMi9i1iijWwDZ1ybEkkuHMltPw-AFNe4EfM8GG-mXlTI4qI6R6RS6rYHluLR5StG_V_r6uOBZxrJ7V3_v2NCO5g0Nl0YTmXEEVcI4PeoXLbWBa8eZ1JvzHSr-JOltHojILhVcZY9RY1OM3unl37gNoz-n0jMobgIQ51LPSv35-tv6SqzCNWINTh6_jvQkmg8UFY1Ce-20hrP00TrN2csxG4Y0YOS4EXmWSwmlObwl8zVA2W2--kf6n05Yi1cJV5JPYLFCFltqI2_ghBjHCYXo3rGAGV81tONRbLDohi7L9HMAHKh1yvdIkCe7EfqzoSXs8cvR7RHlpyZ6NxXlzHj5l0Pd7mDC3SduQx8gRhOv1aP0X3-LtFk5aQyWJk9379KPRcSbvS84CgMDuPjwgRpeMDsdqgbNVWzR7KYrkdJsU_VBK7NBvEtA9LqRmBNSdmd3XEHWBdOUMctVGBLg2OfXtmuCosWSMcg5khPERwG37a4fHSV8sYDyt-zqyKwPhjQ8OyVA5VsjhvwFh6sWjmqe3kAz3HfGq_MLhjRp5OACrlKLLUy6e6eGSplLUeCvx00tzpqIFECAlaXvtyXSYmIUodnboFvUtTbl3dDCwR3eJR0hw3M_mybrOxKtGqxs6Wj50x-zqnXi5TVgqfDMJQ8SILO-QXYgVCCZir314GyaTTlwivP7kCMOnHdmmnOJ6v4Am4-aSLfJ1DrwZLZjhmU3TqEsASkEIbF6JnhA2z8qqjSCH9ZuSFjSWZ5fSaB0xKzN_sIIQEps2y3pm4WjCbExj4Zub_GaOHgvMLy4OQd1JKnals6F8BLSgbxdsIIANMOEVoiWbbkHWpLRlqrYJIfk-PntWIPmk58jhfqoUXtJhY7X-5bz0qtX_-UG4_TT3DUWeieaRnWQl_vDPWMk9dJiiRsyutW0cUkWqNyPYXQSxPWJptd_BHBwyk7jlM_9e5ywQPj_UlQkeYd0ZkbO1so_9jdJwkzgaOoyS5PxnwOCYAr2mIFLJGTgfr8XauCZ10JKd3_I8M3fqzzsF3Y-x2RnSQOb1djAMC4w97BOwVZsb-dsfoDLMZKIBVOOsPywsk26Ebpnvut4n6j7x-0EzixfFP4lgTqPiEw0JdoM9jgeXTwqZxTKW3vRJZkA6Nt9hDeI-IoxHi-Vcgy7p5kiQ7Grp18u8bu7WQsp_4-D49OQ2khZ_laqY_ict_b222sAGAG8zsddZ8WiigWrpG9liTO2WAMfqru4QNFRY1f36T3af33HAPUoaLP4UcBOqKGmiOErJ1nfB62iBdRTi6crUJdyGCjfxXDzZiBkZdpje7QXhjiePgIoVjnD6WXVLsAQYyY4ZUNOxLf4lf2hs9rXQh23zwA8lIWuq63FtY0CYAHLKh5qXcUOm0Zn4qV7xLoEXlCf5ApKON4louv-V1BXLlKqGpnJbd3Ey0OyMokJG_c4taHaK1PMRQvQPafbmBM1eZZ4IW588dkG0IczpIJPmScwLSxCGLlQUGRGwbDiGcdcM78NuojxkHcv6oiMZZgOshVKQAPYgLRpprVF2ZfhyydKDS92g8lZuE.5jeXyyBp3FOfZEgquEFg7g; _T_ANO=c+ZA+9kJ5OSGb+l11ooQ6Nw8NiAnem1ck+7HEbSMoSMHcXEadE6VXxMN4yRtor7+aDJ9cjWlvCuWHMxjo1F4Y5cofU2mI3pYQlYIaHnR2vXH65+cuOGWN7lZQQJK818pwWP0U+rE/Pk8i0X3XB9wxSkOywN9jrZR3ibc75kRxztDkxFkMbeRTs7qtEFOcfuBb9xri9eJhQjoKch3wE+jKcyXdprCmEch2n16C60B4SVo1X5DmiGknOuNbaN9Pr9vbxi5Wd79z8pmgQI15OMP4NICk3xPU/IEpoUf7Y/SPx6IOn+trGjzdWsyCRmcLt8uK7ERjZX4nsy88UFLUtih3Q=="
)

// ============================================================================
// QUERY (dùng payload bạn cung cấp — giữ nguyên các type Long!)
// ============================================================================

const ViewerQuery = `
query viewerInfo($seriesId: Long!, $productId: Long!) {
  viewerInfo(seriesId: $seriesId, productId: $productId) {
    item {
      ...SingleFragment
    }
    seriesItem {
      ...SeriesFragment
    }
    prevItem {
      ...NearItemFragment
    }
    nextItem {
      ...NearItemFragment
    }
    viewerData {
      ...TextViewerData
      ...TalkViewerData
      ...ImageViewerData
      ...VodViewerData
    }
    displayAd {
      ...DisplayAd
    }
  }
}

fragment SingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  series {
    ...SeriesFragment
  }
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  discountRate
  discountRateText
  isShortsDrama
}

fragment SeriesFragment on Series {
  id
  seriesId
  title
  thumbnail
  landThumbnail
  categoryUid
  lang
  category
  categoryType
  subcategoryUid
  subcategory
  badge
  isAllFree
  isWaitfree
  ageGrade
  state
  onIssue
  authors
  description
  pubPeriod
  freeSlideCount
  lastSlideAddedDate
  waitfreeBlockCount
  waitfreePeriodByMinute
  bm
  saleState
  startSaleDt
  saleMethod
  discountRate
  discountRateText
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  translateProperty {
    ...TranslatePropertyFragment
  }
}

fragment ServicePropertyFragment on ServiceProperty {
  viewCount
  readCount
  ratingCount
  ratingSum
  commentCount
  pageContinue {
    ...ContinueInfoFragment
  }
  todayGift {
    ...TodayGift
  }
  preview {
    ...PreviewFragment
  }
  waitfreeTicket {
    ...WaitfreeTicketFragment
  }
  isAlarmOn
  isLikeOn
  ticketCount
  purchasedDate
  lastViewInfo {
    ...LastViewInfoFragment
  }
  purchaseInfo {
    ...PurchaseInfoFragment
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment ContinueInfoFragment on ContinueInfo {
  title
  isFree
  productId
  lastReadProductId
  scheme
  continueProductType
  hasNewSingle
  hasUnreadSingle
}

fragment TodayGift on TodayGift {
  id
  uid
  ticketType
  ticketKind
  ticketCount
  ticketExpireAt
  ticketExpiredText
  isReceived
  seriesId
}

fragment PreviewFragment on Preview {
  item {
    ...PreviewSingleFragment
  }
  nextItem {
    ...PreviewSingleFragment
  }
  usingScroll
}

fragment PreviewSingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
}

fragment OperatorPropertyFragment on OperatorProperty {
  thumbnail
  copy
  helixImpId
  isTextViewer
  selfCensorship
  isBook
  cashInfo {
    discountRate
    setDiscountRate
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment AssetPropertyFragment on AssetProperty {
  bannerImage
  cardImage
  cardTextImage
  cleanImage
  ipxVideo
  bannerSet {
    ...BannerSetFragment
  }
  cardSet {
    ...CardSetFragment
  }
  cardCover {
    ...CardCoverFragment
  }
}

fragment BannerSetFragment on BannerSet {
  backgroundImage
  backgroundColor
  mainImage
  titleImage
}

fragment CardSetFragment on CardSet {
  backgroundColor
  backgroundImage
}

fragment CardCoverFragment on CardCover {
  coverImg
  coverRestricted
}

fragment WaitfreeTicketFragment on WaitfreeTicket {
  chargedPeriod
  chargedCount
  chargedAt
}

fragment LastViewInfoFragment on LastViewInfo {
  isDone
  lastViewDate
  rate
  spineIndex
}

fragment PurchaseInfoFragment on PurchaseInfo {
  purchaseType
  rentExpireDate
  expired
}

fragment TranslatePropertyFragment on TranslateProperty {
  category {
    ...LocaleMapFragment
  }
  sub_category {
    ...LocaleMapFragment
  }
}

fragment LocaleMapFragment on LocaleMap {
  ko
  en
  th
}

fragment NearItemFragment on NearItem {
  productId
  slideType
  ageGrade
  isFree
  title
  thumbnail
}

fragment TextViewerData on TextViewerData {
  type
  atsServerUrl
  metaSecureUrl
  contentsList {
    chapterId
    contentId
    secureUrl
  }
}

fragment TalkViewerData on TalkViewerData {
  type
  talkDownloadData {
    dec
    host
    path
    talkViewerType
  }
}

fragment ImageViewerData on ImageViewerData {
  type
  imageDownloadData {
    ...ImageDownloadData
  }
}

fragment ImageDownloadData on ImageDownloadData {
  files {
    ...ImageDownloadFile
  }
  totalCount
  totalSize
  viewDirection
  gapBetweenImages
  readType
}

fragment ImageDownloadFile on ImageDownloadFile {
  no
  size
  secureUrl
  width
  height
}

fragment VodViewerData on VodViewerData {
  type
  vodDownloadData {
    contentPackId
    drmType
    endpointUrl
    width
    height
    duration
  }
  drmInfo {
    type
    serverType
    error
    fairplayCertificateUrl
    widevineLicenseUrl
    fairplayLicenseUrl
    token
    provider
    assertion
  }
}

fragment DisplayAd on DisplayAd {
  sectionUid
  bannerUid
  treviUid
  momentUid
}
`

// ============================================================================
// MODELS (mình chỉ khai báo những field cần dùng)
// ============================================================================

type GraphQLPayload struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

type GQLResponse struct {
	Data struct {
		ViewerInfo ViewerInfo `json:"viewerInfo"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type ViewerInfo struct {
	Item       Item       `json:"item"`
	SeriesItem *Series    `json:"seriesItem,omitempty"`
	PrevItem   *NearItem  `json:"prevItem,omitempty"`
	NextItem   *NearItem  `json:"nextItem,omitempty"`
	ViewerData ViewerData `json:"viewerData"`
}

type Item struct {
	Title string `json:"title"`
	// thêm field nếu cần về sau
	ProductId int `json:"productId,omitempty"`
	SeriesId  int `json:"seriesId,omitempty"`
}

type Series struct {
	SeriesId int    `json:"seriesId"`
	Title    string `json:"title"`
}

type NearItem struct {
	ProductId int    `json:"productId"`
	Title     string `json:"title"`
	// slideType, etc. omitted — thêm nếu cần
}

type ViewerData struct {
	AtsServerUrl string        `json:"atsServerUrl"`
	MetaSecure   string        `json:"metaSecureUrl,omitempty"`
	ContentsList []ContentPart `json:"contentsList"`
}

type ContentPart struct {
	ChapterId int    `json:"chapterId,omitempty"`
	ContentId int    `json:"contentId,omitempty"`
	SecureUrl string `json:"secureUrl"`
}

type ContentFileResponse struct {
	ContentInfo struct {
		ParagraphList []ParagraphNode `json:"paragraphList"`
	} `json:"contentInfo"`
}

type ParagraphNode struct {
	Type     string          `json:"type"`
	Text     string          `json:"text"`
	Children []ParagraphNode `json:"childParagraphList,omitempty"`
}

// ============================================================================
// MAIN: loop theo nextItem.productId — giữ flow tải/chạy y nguyên code gốc
// ============================================================================

func main() {
	// context tổng cho toàn bộ chạy
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	client := &http.Client{Timeout: 30 * time.Second}
	headers := map[string]string{
		"accept":          "*/*",
		"accept-language": "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7",
		"content-type":    "application/json",
		"origin":          "https://page.kakao.com",
		"referer":         "https://page.kakao.com/",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
		"cookie":          CookieString,
	}

	// Bắt đầu từ đây: set seriesId và productId chương ban đầu bạn muốn crawl
	currentSeries := 61097897
	currentProduct := 65893128

	for {
		fmt.Printf("\n>>> Crawl productId = %d (seriesId=%d)\n", currentProduct, currentSeries)

		// 1) build payload giữ nguyên như bạn yêu cầu
		payload := GraphQLPayload{
			OperationName: "viewerInfo",
			Query:         ViewerQuery,
			Variables: map[string]interface{}{
				"seriesId":  currentSeries,
				"productId": currentProduct,
			},
		}
		payloadBytes, _ := json.Marshal(payload)

		respBody, err := makeRequest(ctx, client, "POST", GraphQLEndpoint, headers, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Printf(">>> ERROR: GraphQL request failed: %v\n", err)
			break
		}

		var gqlResp GQLResponse
		if err := json.Unmarshal(respBody, &gqlResp); err != nil {
			fmt.Printf(">>> ERROR: Parse Error: %v\nBody: %s\n", err, string(respBody))
			break
		}
		if len(gqlResp.Errors) > 0 {
			fmt.Printf(">>> ERROR: GraphQL Error: %s\n", gqlResp.Errors[0].Message)
			break
		}

		viewerData := gqlResp.Data.ViewerInfo.ViewerData
		itemInfo := gqlResp.Data.ViewerInfo.Item
		nextItem := gqlResp.Data.ViewerInfo.NextItem

		if viewerData.AtsServerUrl == "" {
			fmt.Println(">>> ERROR: Không tìm thấy atsServerUrl. Kiểm tra lại Cookie.")
			break
		}

		// Tạo file
		safeTitle := sanitizeFilename(itemInfo.Title)
		fileName := safeTitle + ".txt"

		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf(">>> ERROR: Không thể tạo file %s: %v\n", fileName, err)
			break
		}
		fmt.Printf(">>> Đã tạo file: %s\n", fileName)
		fmt.Printf(">>> Tổng số phần: %d. Bắt đầu tải...\n", len(viewerData.ContentsList))

		// Ghi tiêu đề
		if _, err := file.WriteString(fmt.Sprintf("# %s\n\n", itemInfo.Title)); err != nil {
			fmt.Printf(">>> WARNING: lỗi ghi tiêu đề: %v\n", err)
		}

		// Loop tải các phần
		for i, part := range viewerData.ContentsList {
			if i > 0 {
				time.Sleep(1 * time.Second)
			}
			fullUrl := viewerData.AtsServerUrl + part.SecureUrl
			fmt.Printf(" -> [%d/%d] Đang tải... (Delay 1s)\n", i+1, len(viewerData.ContentsList))

			partBody, err := makeRequest(ctx, client, "GET", fullUrl, headers, nil)
			if err != nil {
				fmt.Printf("   [Lỗi tải phần %d]: %v\n", i+1, err)
				continue
			}

			var contentResp ContentFileResponse
			if err := json.Unmarshal(partBody, &contentResp); err != nil {
				// nếu không parse được JSON -> có thể là image / binary -> bỏ qua
				continue
			}

			for _, node := range contentResp.ContentInfo.ParagraphList {
				text := extractTextRecursive(node)
				text = strings.ReplaceAll(text, "&nbsp;", " ")
				if _, err := file.WriteString(text); err != nil {
					fmt.Printf("   Lỗi ghi file: %v\n", err)
				}
			}
		}

		file.Close()
		fmt.Println(">>> XONG! Kiểm tra file:", fileName)

		// Nếu không có nextItem -> dừng
		if nextItem == nil || nextItem.ProductId == 0 {
			fmt.Println(">>> HẾT CHƯƠNG hoặc không tìm thấy nextItem.")
			break
		}

		// Cập nhật productId để crawl chương tiếp theo
		fmt.Printf(">>> Next: %d (%s)\n", nextItem.ProductId, nextItem.Title)
		currentProduct = nextItem.ProductId

		// delay an toàn trước lần tiếp theo
		time.Sleep(2 * time.Second)
	}
}

// ============================================================================
// HELPERS
// ============================================================================

func makeRequest(ctx context.Context, client *http.Client, method, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status %d | Body: %s", resp.StatusCode, string(respBytes))
	}
	return respBytes, nil
}

func extractTextRecursive(node ParagraphNode) string {
	var sb strings.Builder
	if node.Text != "" {
		sb.WriteString(node.Text)
	}
	if node.Type == "BR" {
		sb.WriteString("\n")
	}
	for _, child := range node.Children {
		sb.WriteString(extractTextRecursive(child))
	}
	switch node.Type {
	case "P", "DIV":
		sb.WriteString("\n")
	case "H3":
		sb.WriteString("\n\n### ")
	}
	return sb.String()
}

func sanitizeFilename(name string) string {
	reg, _ := regexp.Compile(`[\\/:*?"<>|]`)
	safe := reg.ReplaceAllString(name, "_")
	return strings.TrimSpace(safe)
}
